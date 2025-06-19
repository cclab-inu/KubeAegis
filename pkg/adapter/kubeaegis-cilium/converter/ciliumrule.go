package converter

import (
	"fmt"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cilium/cilium/pkg/policy/api"
)

// getIngressDeny generates ingress deny rules from KubeAegisPolicy IntentRequests
func getIngressDeny(intentRequest v1.IntentRequest) ([]api.IngressDenyRule, error) {
	var ingressDenyRules []api.IngressDenyRule

	for _, from := range intentRequest.Rule.From {
		ingressRule := api.IngressDenyRule{}

		switch from.Kind {
		case "endpoint":
			if len(from.Labels) > 0 {
				// Convert MatchLabels to Cilium EndpointSelector
				selector := api.NewESFromMatchRequirements(from.Labels, nil)
				ingressRule.FromEndpoints = append(ingressRule.FromEndpoints, selector)
			}
		case "entities":
			if len(from.Args) > 0 {
				// Add each entity to the FromEntities slice
				for _, entity := range from.Args {
					ingressRule.FromEntities = append(ingressRule.FromEntities, api.Entity(entity))
				}
			}
		case "cidr":
			if len(from.Args) > 0 {
				// Add each CIDR to the FromCIDR slice
				for _, cidr := range from.Args {
					ingressRule.FromCIDR = append(ingressRule.FromCIDR, api.CIDR(cidr))
				}
			}

		default:
			return nil, fmt.Errorf("unsupported kind: %s", from.Kind)
		}

		ingressDenyRules = append(ingressDenyRules, ingressRule)
	}
	return ingressDenyRules, nil
}

func getIngress(intentRequest v1.IntentRequest) ([]api.IngressRule, error) {
	var ingressRules []api.IngressRule
	ingressRule := api.IngressRule{}

	for _, from := range intentRequest.Rule.From {
		switch from.Kind {
		case "endpoint":
			if len(from.Labels) > 0 {
				// Convert MatchLabels to Cilium EndpointSelector
				selector := api.NewESFromMatchRequirements(from.Labels, nil)
				ingressRule.FromEndpoints = append(ingressRule.FromEndpoints, selector)
			}
		default:
			return nil, fmt.Errorf("unsupported kind: %s", from.Kind)
		}
		ingressRules = append(ingressRules, ingressRule)
	}
	return ingressRules, nil
}

func getEgressDeny(intentRequest v1.IntentRequest) ([]api.EgressDenyRule, error) {
	var egressDenyRules []api.EgressDenyRule
	for _, to := range intentRequest.Rule.To {
		switch to.Kind {
		case "port":
			// Convert to PortProtocol format for cilium
			portProtocols := make([]api.PortProtocol, 0)
			portProtocol := api.PortProtocol{
				Port:     to.Port,
				Protocol: api.L4Proto(to.Protocol),
			}
			portProtocols = append(portProtocols, portProtocol)

			egressDenyRule := api.EgressDenyRule{
				ToPorts: []api.PortDenyRule{
					{
						Ports: portProtocols,
					},
				},
			}
			egressDenyRules = append(egressDenyRules, egressDenyRule)
		}
	}
	return egressDenyRules, nil
}
