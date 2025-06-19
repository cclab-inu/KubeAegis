package converter

import (
	"fmt"
	"strconv"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	calico "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	crdnum "github.com/projectcalico/api/pkg/lib/numorstring"
)

// getIngressRules generates ingress rules from KubeAegisPolicy IntentRequests
func getIngressRules(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	var ingressRules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		ingressRule := calico.Rule{}

		switch from.Kind {
		case "pod":
			if len(from.Labels) > 0 {
				ingressRule.Source.Selector = createSelector(from.Labels)
			}
		case "namespace", "serviceAccounts":
			if len(from.Labels) > 0 {
				ingressRule.Source.NamespaceSelector = createSelector(from.Labels)
			}
		case "cidr":
			if len(from.Args) > 0 {
				ingressRule.Source.Nets = from.Args
			}
		case "protocol":
			if from.Protocol != "" {
				p := crdnum.ProtocolFromString(from.Protocol)
				ingressRule.Protocol = &p
			}
		default:
			return nil, fmt.Errorf("unsupported kind: %s", from.Kind)
		}

		ingressRules = append(ingressRules, ingressRule)
	}
	return ingressRules, nil
}

// getEgressRules generates egress rules from KubeAegisPolicy IntentRequests
func getEgressRules(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	var egressRules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		egressRule := calico.Rule{}

		switch to.Kind {
		case "pod":
			if len(to.Labels) > 0 {
				egressRule.Destination.Selector = createSelector(to.Labels)
			}
		case "namespace", "serviceAccounts":
			if len(to.Labels) > 0 {
				egressRule.Destination.NamespaceSelector = createSelector(to.Labels)
			}
		case "cidr":
			if len(to.Args) > 0 {
				egressRule.Destination.Nets = to.Args
			}
		case "port":
			if to.Port != "" {
				pi, err := strconv.Atoi(to.Port)
				if err != nil {
					return nil, fmt.Errorf("invalid port %q: %v", to.Port, err)
				}
				pp := crdnum.SinglePort(uint16(pi))
				egressRule.Destination.Ports = []crdnum.Port{pp}

			}
			if to.Protocol != "" {
				protocol := crdnum.ProtocolFromString(to.Protocol)
				egressRule.Protocol = &protocol
			}
		default:
			return nil, fmt.Errorf("unsupported kind: %s", to.Kind)
		}

		egressRules = append(egressRules, egressRule)
	}
	return egressRules, nil
}

// getHTTPRules generates HTTP rules from ActionPoints in KubeAegisPolicy
func getHTTPRules(actionPoints []v1.ActionPoint) ([]calico.HTTPMatch, error) {
	var httpRules []calico.HTTPMatch

	for _, ap := range actionPoints {
		if ap.SubType == "http" {
			httpRule := calico.HTTPMatch{
				Methods: ap.Resource.Methods,
				Paths:   createHTTPPaths(ap.Resource.Path),
			}
			httpRules = append(httpRules, httpRule)
		}
	}
	return httpRules, nil
}

// createHTTPPaths converts path strings to HTTPPath objects
func createHTTPPaths(paths []string) []calico.HTTPPath {
	var httpPaths []calico.HTTPPath
	for _, path := range paths {
		httpPaths = append(httpPaths, calico.HTTPPath{Exact: path})
	}
	return httpPaths
}

// createSelector generates a selector string from given labels
func createSelector(labels map[string]string) string {
	var selector string
	for key, value := range labels {
		if selector != "" {
			selector += " && "
		}
		selector += fmt.Sprintf("%s == '%s'", key, value)
	}
	return selector
}
