package recommandpool

import (
	"strconv"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	ciliumapi "github.com/cilium/cilium/pkg/policy/api"
	calico "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	calicoapilib "github.com/projectcalico/api/pkg/lib/numorstring"
	"github.com/projectcalico/libcalico-go/lib/numorstring"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var ruleDescription string

// ================================================================================================================
// ----------------------------
// Cilium
// ----------------------------
// ================================================================================================================

// ----------------------------
// Network Policy: Endpoint
// ----------------------------

// CreateIngressEndpointSelectorRule generates an ingress rule with endpoint selectors for Cilium
func CreateIngressEndpointSelectorRule(intentRequest v1.IntentRequest) ([]ciliumapi.IngressRule, error) {
	ruleDescription = "This function generates ingress rules with endpoint selectors for Cilium. It processes the 'From' field in the intent request, which specifies the source endpoints that need to be allowed access. The function extracts labels from the 'From' field to create endpoint selectors, which are then used to define ingress rules that allow traffic from the specified endpoints."
	var ingressRules []ciliumapi.IngressRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "endpoint" && len(from.Labels) > 0 {
			selector := ciliumapi.NewESFromMatchRequirements(from.Labels, nil)
			Rule := ciliumapi.IngressRule{
				IngressCommonRule: ciliumapi.IngressCommonRule{
					FromEndpoints: []ciliumapi.EndpointSelector{selector},
				},
			}
			ingressRules = append(ingressRules, Rule)
		}
	}
	return ingressRules, nil
}

// CreateIngressDenyEndpointSelectorRule generates an ingress deny rule with endpoint selectors for Cilium
func CreateIngressDenyEndpointSelectorRule(intentRequest v1.IntentRequest) ([]ciliumapi.IngressDenyRule, error) {
	ruleDescription = "This function generates ingress deny rules with endpoint selectors for Cilium. It processes the 'From' field in the intent request, which specifies the source endpoints that need to be denied access. The function extracts labels from the 'From' field to create endpoint selectors, which are then used to define ingress deny rules that prevent traffic from the specified endpoints."
	var ingressDenyRules []ciliumapi.IngressDenyRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "endpoint" && len(from.Labels) > 0 {
			selector := ciliumapi.NewESFromMatchRequirements(from.Labels, nil)
			Rule := ciliumapi.IngressDenyRule{
				IngressCommonRule: ciliumapi.IngressCommonRule{
					FromEndpoints: []ciliumapi.EndpointSelector{selector},
				},
			}
			ingressDenyRules = append(ingressDenyRules, Rule)
		}
	}
	return ingressDenyRules, nil
}

// CreateEgressEndpointSelectorRule generates an egress rule with endpoint selectors for Cilium
func CreateEgressEndpointSelectorRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressRule, error) {
	ruleDescription = "This function generates egress rules with endpoint selectors for Cilium. It processes the 'To' field in the intent request, which specifies the destination endpoints that need to be allowed access. The function extracts labels from the 'To' field to create endpoint selectors, which are then used to define egress rules that allow traffic to the specified endpoints."
	var egressRules []ciliumapi.EgressRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "endpoint" && len(to.Labels) > 0 {
			selector := ciliumapi.NewESFromMatchRequirements(to.Labels, nil)
			Rule := ciliumapi.EgressRule{
				EgressCommonRule: ciliumapi.EgressCommonRule{
					ToEndpoints: []ciliumapi.EndpointSelector{selector},
				},
			}
			egressRules = append(egressRules, Rule)
		}
	}
	return egressRules, nil
}

// CreateEgressDenyEndpointSelectorRule generates an egress deny rule with endpoint selectors for Cilium
func CreateEgressDenyEndpointSelectorRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressDenyRule, error) {
	ruleDescription = "This function generates egress deny rules with endpoint selectors for Cilium. It processes the 'To' field in the intent request, which specifies the destination endpoints that need to be denied access. The function extracts labels from the 'To' field to create endpoint selectors, which are then used to define egress deny rules that prevent traffic to the specified endpoints."
	var egressDenyRules []ciliumapi.EgressDenyRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "endpoint" && len(to.Labels) > 0 {
			selector := ciliumapi.NewESFromMatchRequirements(to.Labels, nil)
			Rule := ciliumapi.EgressDenyRule{
				EgressCommonRule: ciliumapi.EgressCommonRule{
					ToEndpoints: []ciliumapi.EndpointSelector{selector},
				},
			}
			egressDenyRules = append(egressDenyRules, Rule)
		}
	}
	return egressDenyRules, nil
}

// ----------------------------
// Network Policy: Entities
// ----------------------------

// CreateIngressEntitiesRule generates an ingress rule with entities for Cilium
func CreateIngressEntitiesRule(intentRequest v1.IntentRequest) ([]ciliumapi.IngressRule, error) {
	ruleDescription = "This function generates ingress rules with entities for Cilium. It processes the 'From' field in the intent request, which specifies the source entities that need to be allowed access. The function iterates over the entities listed in the 'From' field and creates ingress rules to allow traffic from these entities."
	var ingressRules []ciliumapi.IngressRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "entities" && len(from.Args) > 0 {
			Rule := ciliumapi.IngressRule{}
			for _, entity := range from.Args {
				Rule.FromEntities = append(Rule.FromEntities, ciliumapi.Entity(entity))
			}
			ingressRules = append(ingressRules, Rule)
		}
	}
	return ingressRules, nil
}

// CreateIngressDenyEntitiesRule generates an ingress deny rule with entities for Cilium
func CreateIngressDenyEntitiesRule(intentRequest v1.IntentRequest) ([]ciliumapi.IngressDenyRule, error) {
	ruleDescription = "This function generates ingress deny rules with entities for Cilium. It processes the 'From' field in the intent request, which specifies the source entities that need to be denied access. The function iterates over the entities listed in the 'From' field and creates ingress deny rules to block traffic from these entities."
	var ingressDenyRules []ciliumapi.IngressDenyRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "entities" && len(from.Args) > 0 {
			Rule := ciliumapi.IngressDenyRule{}
			for _, entity := range from.Args {
				Rule.FromEntities = append(Rule.FromEntities, ciliumapi.Entity(entity))
			}
			ingressDenyRules = append(ingressDenyRules, Rule)
		}
	}
	return ingressDenyRules, nil
}

// CreateEgressEntitiesRule generates an egress rule with entities for Cilium
func CreateEgressEntitiesRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressRule, error) {
	ruleDescription = "This function generates egress rules with entities for Cilium. It processes the 'To' field in the intent request, which specifies the destination entities that need to be allowed access. The function iterates over the entities listed in the 'To' field and creates egress rules to allow traffic to these entities."
	var egressRules []ciliumapi.EgressRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "entities" && len(to.Args) > 0 {
			Rule := ciliumapi.EgressRule{}
			for _, entity := range to.Args {
				Rule.ToEntities = append(Rule.ToEntities, ciliumapi.Entity(entity))
			}
			egressRules = append(egressRules, Rule)
		}
	}
	return egressRules, nil
}

// CreateEgressDenyEntitiesRule generates an egress deny rule with entities for Cilium
func CreateEgressDenyEntitiesRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressDenyRule, error) {
	ruleDescription = "This function generates egress deny rules with entities for Cilium. It processes the 'To' field in the intent request, which specifies the destination entities that need to be denied access. The function iterates over the entities listed in the 'To' field and creates egress deny rules to block traffic to these entities."
	var egressDenyRules []ciliumapi.EgressDenyRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "entities" && len(to.Args) > 0 {
			Rule := ciliumapi.EgressDenyRule{}
			for _, entity := range to.Args {
				Rule.ToEntities = append(Rule.ToEntities, ciliumapi.Entity(entity))
			}
			egressDenyRules = append(egressDenyRules, Rule)
		}
	}
	return egressDenyRules, nil
}

// ----------------------------
// Network Policy: CIDR
// ----------------------------

// CreateIngressCIDRRule generates an ingress rule with CIDR blocks for Cilium
func CreateIngressCIDRRule(intentRequest v1.IntentRequest) ([]ciliumapi.IngressRule, error) {
	ruleDescription = "This function generates ingress rules with CIDR blocks for Cilium. It processes the 'From' field in the intent request, which specifies the source CIDR blocks that need to be allowed access. The function iterates over the CIDR blocks listed in the 'From' field and creates ingress rules to allow traffic from these IP ranges."
	var ingressRules []ciliumapi.IngressRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "cidr" && len(from.Args) > 0 {
			Rule := ciliumapi.IngressRule{}
			for _, cidr := range from.Args {
				Rule.FromCIDR = append(Rule.FromCIDR, ciliumapi.CIDR(cidr))
			}
			ingressRules = append(ingressRules, Rule)
		}
	}
	return ingressRules, nil
}

// CreateIngressDenyCIDRRule generates an ingress deny rule with CIDR blocks for Cilium
func CreateIngressDenyCIDRRule(intentRequest v1.IntentRequest) ([]ciliumapi.IngressDenyRule, error) {
	ruleDescription = "This function generates ingress deny rules with CIDR blocks for Cilium. It processes the 'From' field in the intent request, which specifies the source CIDR blocks that need to be denied access. The function iterates over the CIDR blocks listed in the 'From' field and creates ingress deny rules to block traffic from these IP ranges."
	var ingressDenyRules []ciliumapi.IngressDenyRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "cidr" && len(from.Args) > 0 {
			Rule := ciliumapi.IngressDenyRule{}
			for _, cidr := range from.Args {
				Rule.FromCIDR = append(Rule.FromCIDR, ciliumapi.CIDR(cidr))
			}
			ingressDenyRules = append(ingressDenyRules, Rule)
		}
	}
	return ingressDenyRules, nil
}

// CreateEgressCIDRRule generates an egress rule with CIDR blocks for Cilium
func CreateEgressCIDRRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressRule, error) {
	ruleDescription = "This function generates egress rules with CIDR blocks for Cilium. It processes the 'To' field in the intent request, which specifies the destination CIDR blocks that need to be allowed access. The function iterates over the CIDR blocks listed in the 'To' field and creates egress rules to allow traffic to these IP ranges."
	var egressRules []ciliumapi.EgressRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "cidr" && len(to.Args) > 0 {
			Rule := ciliumapi.EgressRule{}
			for _, cidr := range to.Args {
				Rule.ToCIDR = append(Rule.ToCIDR, ciliumapi.CIDR(cidr))
			}
			egressRules = append(egressRules, Rule)
		}
	}
	return egressRules, nil
}

// CreateEgressDenyCIDRRule generates an egress deny rule with CIDR blocks for Cilium
func CreateEgressDenyCIDRRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressDenyRule, error) {
	ruleDescription = "This function generates egress deny rules with CIDR blocks for Cilium. It processes the 'To' field in the intent request, which specifies the destination CIDR blocks that need to be denied access. The function iterates over the CIDR blocks listed in the 'To' field and creates egress deny rules to block traffic to these IP ranges."
	var egressDenyRules []ciliumapi.EgressDenyRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "cidr" && len(to.Args) > 0 {
			Rule := ciliumapi.EgressDenyRule{}
			for _, cidr := range to.Args {
				Rule.ToCIDR = append(Rule.ToCIDR, ciliumapi.CIDR(cidr))
			}
			egressDenyRules = append(egressDenyRules, Rule)
		}
	}
	return egressDenyRules, nil
}

// ----------------------------
// Network Policy: FQDNs
// ----------------------------

// CreateEgressFQDNsRule generates an egress rule with FQDNs for Cilium
func CreateEgressFQDNsRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressRule, error) {
	ruleDescription = "This function generates egress rules with FQDNs for Cilium. It processes the 'To' field in the intent request, which specifies the destination FQDNs that need to be allowed access. The function iterates over the FQDNs listed in the 'To' field and creates egress rules to allow traffic to these FQDNs."
	var egressRules []ciliumapi.EgressRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "fqdn" && len(to.Args) > 0 {
			var fqdnRules []ciliumapi.FQDNSelector
			for _, fqdn := range to.Args {
				fqdnRules = append(fqdnRules, ciliumapi.FQDNSelector{
					MatchName: fqdn,
				})
			}
			Rule := ciliumapi.EgressRule{
				ToFQDNs: fqdnRules,
			}
			egressRules = append(egressRules, Rule)
		}
	}
	return egressRules, nil
}

// ----------------------------
// Network Policy: PORT
// ----------------------------

// CreateEgressDenyPortRule generates an egress deny rule with specific ports for Cilium
func CreateEgressDenyPortRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressDenyRule, error) {
	ruleDescription = "This function generates egress deny rules with specific ports for Cilium. It processes the 'To' field in the intent request, which specifies the destination ports that need to be denied access. The function creates egress deny rules that block traffic to the specified ports, ensuring that traffic cannot reach these destinations."
	var egressDenyRules []ciliumapi.EgressDenyRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "port" {
			portProtocols := []ciliumapi.PortProtocol{
				{
					Port:     to.Port,
					Protocol: ciliumapi.L4Proto(to.Protocol),
				},
			}
			Rule := ciliumapi.EgressDenyRule{
				ToPorts: []ciliumapi.PortDenyRule{
					{
						Ports: portProtocols,
					},
				},
			}
			egressDenyRules = append(egressDenyRules, Rule)
		}
	}
	return egressDenyRules, nil
}

// ----------------------------
// Network Policy: ICMP
// ----------------------------

// CreateEgressICMPRule generates an egress rule with ICMP types for Cilium
func CreateEgressICMPField(intentRequest v1.IntentRequest) ([]ciliumapi.EgressRule, error) {
	ruleDescription = "This function generates egress rules with ICMP types for Cilium. It processes the 'To' field in the intent request, which specifies the destination ICMP types that need to be allowed access. The function iterates over the ICMP types listed in the 'To' field and creates egress rules to allow traffic to these ICMP types."
	var egressRules []ciliumapi.EgressRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "icmp" && len(to.Args) > 0 {
			var icmps []ciliumapi.ICMPRule
			for _, icmp := range to.Args {
				icmps = append(icmps, ciliumapi.ICMPRule{
					Fields: []ciliumapi.ICMPField{
						{
							Family: icmp, // Default family, modify as needed
							Type: func() *intstr.IntOrString {
								val := intstr.FromInt(0) // 숫자 타입일 경우
								return &val
							}(),
						},
					},
				})
			}
			Rule := ciliumapi.EgressRule{
				ICMPs: icmps,
			}
			egressRules = append(egressRules, Rule)
		}
	}
	return egressRules, nil
}

// CreateEgressDenyICMPRule generates an egress deny rule with ICMP types for Cilium
func CreateEgressDenyICMPField(intentRequest v1.IntentRequest) ([]ciliumapi.EgressDenyRule, error) {
	ruleDescription = "This function generates egress deny rules with ICMP types for Cilium. It processes the 'To' field in the intent request, which specifies the destination ICMP types that need to be denied access. The function iterates over the ICMP types listed in the 'To' field and creates egress deny rules to block traffic to these ICMP types."
	var egressDenyRules []ciliumapi.EgressDenyRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "icmp" && len(to.Args) > 0 {
			var icmps []ciliumapi.ICMPRule
			for _, icmp := range to.Args {
				icmps = append(icmps, ciliumapi.ICMPRule{
					Fields: []ciliumapi.ICMPField{
						{
							Family: icmp, // Default family, modify as needed
							Type: func() *intstr.IntOrString {
								val := intstr.FromInt(0) // 숫자 타입일 경우
								return &val
							}(),
						},
					},
				})
			}
			Rule := ciliumapi.EgressDenyRule{
				ICMPs: icmps,
			}
			egressDenyRules = append(egressDenyRules, Rule)
		}
	}
	return egressDenyRules, nil
}

// CreateIngressICMPField generates an ingress rule with ICMP types for Cilium
func CreateIngressICMPField(intentRequest v1.IntentRequest) ([]ciliumapi.IngressRule, error) {
	ruleDescription = "This function generates ingress rules with ICMP types for Cilium. It processes the 'From' field in the intent request, which specifies the source ICMP types that need to be allowed access. The function iterates over the ICMP types listed in the 'From' field and creates ingress rules to allow traffic from these ICMP types."
	var ingressRules []ciliumapi.IngressRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "icmps" && len(from.Args) > 0 {
			var icmps []ciliumapi.ICMPRule
			for _, icmp := range from.Args {
				icmps = append(icmps, ciliumapi.ICMPRule{
					Fields: []ciliumapi.ICMPField{
						{
							Family: icmp, // Default family, modify as needed
							Type: func() *intstr.IntOrString {
								val := intstr.FromInt(0) // 숫자 타입일 경우
								return &val
							}(),
						},
					},
				})
			}
			Rule := ciliumapi.IngressRule{
				ICMPs: icmps,
			}
			ingressRules = append(ingressRules, Rule)
		}
	}
	return ingressRules, nil
}

// CreateIngressDenyICMPField generates an ingress deny rule with ICMP types for Cilium
func CreateIngressDenyICMPField(intentRequest v1.IntentRequest) ([]ciliumapi.IngressDenyRule, error) {
	ruleDescription = "This function generates ingress deny rules with ICMP types for Cilium. It processes the 'From' field in the intent request, which specifies the source ICMP types that need to be denied access. The function iterates over the ICMP types listed in the 'From' field and creates ingress deny rules to block traffic from these ICMP types."
	var ingressDenyRules []ciliumapi.IngressDenyRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "icmps" && len(from.Args) > 0 {
			var icmps []ciliumapi.ICMPRule
			for _, icmp := range from.Args {
				icmps = append(icmps, ciliumapi.ICMPRule{
					Fields: []ciliumapi.ICMPField{
						{
							Family: icmp, // Default family, modify as needed
							Type: func() *intstr.IntOrString {
								val := intstr.FromInt(0) // 숫자 타입일 경우
								return &val
							}(),
						},
					},
				})
			}
			Rule := ciliumapi.IngressDenyRule{
				ICMPs: icmps,
			}
			ingressDenyRules = append(ingressDenyRules, Rule)
		}
	}
	return ingressDenyRules, nil
}

// ----------------------------
// Network Policy: Services
// ----------------------------

// CreateEgressServiceRule generates an egress rule with services for Cilium
func CreateEgressServiceRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressDenyRule, error) {
	ruleDescription = "This function generates egress deny rules with services for Cilium. It processes the 'To' field in the intent request, which specifies the destination services that need to be denied access. The function iterates over the services listed in the 'To' field and creates egress deny rules to block traffic to these services."
	var egressDenyRules []ciliumapi.EgressDenyRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "service" && len(to.Args) > 0 {
			Rule := ciliumapi.EgressDenyRule{
				EgressCommonRule: ciliumapi.EgressCommonRule{
					ToServices: []ciliumapi.Service{
						{
							K8sService: &ciliumapi.K8sServiceNamespace{
								ServiceName: to.Args[0],
								Namespace:   to.Args[1],
							},
						},
					},
				},
			}
			egressDenyRules = append(egressDenyRules, Rule)
		}
	}
	return egressDenyRules, nil
}

// ----------------------------
// Network Policy: DNS
// ----------------------------

// CreateEgressDNSRule generates an egress rule with DNS constraints for Cilium
func CreateEgressDNSRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressRule, error) {
	ruleDescription = "This function generates egress rules with DNS constraints for Cilium. It processes the 'To' field in the intent request, which specifies the destination DNS constraints. The function creates egress rules that block traffic based on these DNS constraints."
	var egressRules []ciliumapi.EgressRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "dns" && len(to.Args) > 0 {
			dnsRules := []ciliumapi.PortRuleDNS{}
			for _, dns := range to.Args {
				dnsRules = append(dnsRules, ciliumapi.PortRuleDNS{
					MatchName: dns,
				})
			}
			Rule := ciliumapi.EgressRule{
				ToPorts: []ciliumapi.PortRule{
					{
						Rules: &ciliumapi.L7Rules{
							DNS: dnsRules,
						},
					},
				},
			}
			egressRules = append(egressRules, Rule)
		}
	}
	return egressRules, nil
}

// CreateIngressDNSRule generates an ingress rule with DNS constraints for Cilium
func CreateIngressDNSRule(intentRequest v1.IntentRequest) ([]ciliumapi.IngressRule, error) {
	ruleDescription = "This function generates ingress rules with DNS constraints for Cilium. It processes the 'From' field in the intent request, which specifies the source DNS constraints. The function creates ingress rules that allow traffic based on these DNS constraints."
	var ingressRules []ciliumapi.IngressRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "dns" && len(from.Args) > 0 {
			dnsRules := []ciliumapi.PortRuleDNS{}
			for _, dns := range from.Args {
				dnsRules = append(dnsRules, ciliumapi.PortRuleDNS{
					MatchName: dns,
				})
			}
			Rule := ciliumapi.IngressRule{
				ToPorts: []ciliumapi.PortRule{
					{
						Rules: &ciliumapi.L7Rules{
							DNS: dnsRules,
						},
					},
				},
			}
			ingressRules = append(ingressRules, Rule)
		}
	}
	return ingressRules, nil
}

// ----------------------------
// Network Policy: HTTPS
// ----------------------------

// CreateEgressHTTPSRule generates an egress rule with HTTPS constraints for Cilium
func CreateEgressHTTPSRule(intentRequest v1.IntentRequest) ([]ciliumapi.EgressRule, error) {
	ruleDescription = "This function generates egress rules with HTTPS constraints for Cilium. It processes the 'To' field in the intent request, which specifies the destination HTTPS constraints. The function creates egress rules that block traffic based on these HTTPS constraints."
	var egressRules []ciliumapi.EgressRule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "https" && len(to.Args) > 0 {
			httpRules := []ciliumapi.PortRuleHTTP{}
			for _, https := range to.Args {
				httpRules = append(httpRules, ciliumapi.PortRuleHTTP{
					Method: https,
				})
			}
			Rule := ciliumapi.EgressRule{
				ToPorts: []ciliumapi.PortRule{
					{
						Rules: &ciliumapi.L7Rules{
							HTTP: httpRules,
						},
					},
				},
			}
			egressRules = append(egressRules, Rule)
		}
	}
	return egressRules, nil
}

// CreateIngressHTTPSRule generates an ingress rule with HTTPS constraints for Cilium
func CreateIngressHTTPSRule(intentRequest v1.IntentRequest) ([]ciliumapi.IngressRule, error) {
	ruleDescription = "This function generates ingress rules with HTTPS constraints for Cilium. It processes the 'From' field in the intent request, which specifies the source HTTPS constraints. The function creates ingress rules that allow traffic based on these HTTPS constraints."
	var ingressRules []ciliumapi.IngressRule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "https" && len(from.Args) > 0 {
			httpRules := []ciliumapi.PortRuleHTTP{}
			for _, https := range from.Args {
				httpRules = append(httpRules, ciliumapi.PortRuleHTTP{
					Method: https,
				})
			}
			Rule := ciliumapi.IngressRule{
				ToPorts: []ciliumapi.PortRule{
					{
						Rules: &ciliumapi.L7Rules{
							HTTP: httpRules,
						},
					},
				},
			}
			ingressRules = append(ingressRules, Rule)
		}
	}
	return ingressRules, nil
}

// ================================================================================================================
// ----------------------------
// Calico (21)
// ----------------------------
// ================================================================================================================

// ----------------------------
// Calico Policy: Pod Selector
// ----------------------------

// CreateIngressPodSelectorRule generates an ingress rule with pod selectors
func CreateIngressPodSelectorRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules with pod selectors for Calico. It processes the 'From' field in the intent request, which specifies the source pods that need to be allowed access. The function extracts labels from the 'From' field to create pod selectors, which are then used to define ingress rules that permit traffic from the specified pods. This allows for fine-grained control over which pods are permitted to send traffic to the selected pods."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "pod" && len(from.Labels) > 0 {
			Rule := calico.Rule{
				Action: "Allow",
				Source: calico.EntityRule{
					Selector: ExtractGiveFormatsSelector(from.Labels),
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateEgressPodSelectorRule generates an egress rule with pod selectors
func CreateEgressPodSelectorRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules with pod selectors for Calico. It processes the 'To' field in the intent request, which specifies the destination pods that need to be allowed access. The function extracts labels from the 'To' field to create pod selectors, which are then used to define egress rules that permit traffic to the specified pods. This ensures that only traffic destined for specific pods is allowed, providing targeted egress control."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "pod" && len(to.Labels) > 0 {
			Rule := calico.Rule{
				Action: "Allow",
				Destination: calico.EntityRule{
					Selector: ExtractGiveFormatsSelector(to.Labels),
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// ----------------------------
// Calico Policy: Namespace Selector
// ----------------------------

// CreateIngressNamespaceSelectorRule generates an ingress rule with namespace selectors
func CreateIngressNamespaceSelectorRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules with namespace selectors for Calico. It processes the 'From' field in the intent request, which specifies the source namespaces that need to be allowed access. The function extracts labels from the 'From' field to create namespace selectors, which are then used to define ingress rules that permit traffic from the specified namespaces. This is useful for controlling traffic based on the namespace of the source pods."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "namespace" && len(from.Labels) > 0 {
			Rule := calico.Rule{
				Action: "Allow",
				Source: calico.EntityRule{
					NamespaceSelector: ExtractGiveFormatsSelector(from.Labels),
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateEgressNamespaceSelectorRule generates an egress rule with namespace selectors
func CreateEgressNamespaceSelectorRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules with namespace selectors for Calico. It processes the 'To' field in the intent request, which specifies the destination namespaces that need to be allowed access. The function extracts labels from the 'To' field to create namespace selectors, which are then used to define egress rules that permit traffic to the specified namespaces. This helps in managing egress traffic based on the destination namespace."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "namespace" && len(to.Labels) > 0 {
			Rule := calico.Rule{
				Action: "Allow",
				Destination: calico.EntityRule{
					NamespaceSelector: ExtractGiveFormatsSelector(to.Labels),
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// ----------------------------
// Calico Policy: ServiceAccount
// ----------------------------

// CreateIngressServiceAccountsSelectorRule generates an ingress rule with ServiceAccount selectors
func CreateSourceServiceAccountsRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules with ServiceAccount selectors for Calico. It processes the 'From' field in the intent request, which specifies the source service accounts that need to be allowed access. The function extracts labels from the 'From' field to create service account selectors, which are then used to define ingress rules that permit traffic from the specified service accounts. This is crucial for controlling access based on service accounts, ensuring that only pods running as specific service accounts can send traffic."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "serviceAccounts" && len(from.Labels) > 0 {
			Rule := calico.Rule{
				Action: "Allow",
				Source: calico.EntityRule{
					ServiceAccounts: &calico.ServiceAccountMatch{
						Names:    []string{},
						Selector: ExtractGiveFormatsSelector(from.Labels),
					},
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateDestinationServiceAccountsRule generates an egress rule with ServiceAccount selectors
func CreateDestinationServiceAccountsRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules with ServiceAccount selectors for Calico. It processes the 'To' field in the intent request, which specifies the destination service accounts that need to be allowed access. The function extracts labels from the 'To' field to create service account selectors, which are then used to define egress rules that permit traffic to the specified service accounts. This ensures that only traffic destined for pods running as specific service accounts is allowed."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "serviceAccounts" && len(to.Labels) > 0 {
			Rule := calico.Rule{
				Action: "Allow",
				Destination: calico.EntityRule{
					ServiceAccounts: &calico.ServiceAccountMatch{
						Names:    []string{},
						Selector: ExtractGiveFormatsSelector(to.Labels),
					},
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// ----------------------------
// Calico Policy: Protocol
// ----------------------------

// CreateIngressProtocolRule generates an ingress rule with specific protocol
func CreateIngressProtocolRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules with specific IP protocol for Calico. It processes the 'From' field in the intent request, which specifies the source protocols that need to be allowed access. The function extracts the protocol details from the 'From' field and creates ingress rules that permit traffic from the specified protocols. This is useful for enforcing protocol-specific ingress policies. Must be one of these string values: \"TCP\", \"UDP\", \"ICMP\", \"ICMPv6\", \"SCTP\", \"UDPLite\" or an integer in the range 1-255."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "protocol" && from.Protocol != "" {
			protocol := numorstring.ProtocolFromString(from.Protocol)
			calicoProtocol := convertProtocol(protocol)
			Rule := calico.Rule{
				Protocol: &calicoProtocol,
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

func CreateIngressNotProtocolRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules that deny traffic from specified protocols for Calico. It processes the 'From' field in the intent request, which lists the protocols that should be denied. The function extracts protocol details from the 'From' field and creates rules to block ingress traffic using these protocols, enhancing security by excluding unwanted protocol traffic."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "protocol" && from.Protocol != "" {
			protocol := numorstring.ProtocolFromString(from.Protocol)
			calicoProtocol := convertProtocol(protocol)
			Rule := calico.Rule{
				NotProtocol: &calicoProtocol,
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateEgressProtocolRule generates an egress rule with specific protocol
func CreateEgressProtocolRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules with specific IP protocol for Calico. It processes the 'To' field in the intent request, which specifies the destination protocols that need to be allowed access. The function extracts the protocol details from the 'To' field and creates egress rules that permit traffic to the specified protocols. This helps in controlling egress traffic based on the protocol used. Must be one of these string values: \"TCP\", \"UDP\", \"ICMP\", \"ICMPv6\", \"SCTP\", \"UDPLite\" or an integer in the range 1-255."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "protocol" && to.Protocol != "" {
			protocol := numorstring.ProtocolFromString(to.Protocol)
			calicoProtocol := convertProtocol(protocol)
			Rule := calico.Rule{
				Action:   "Allow",
				Protocol: &calicoProtocol,
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateEgressProtocolRule generates an egress rule with specific protocol
func CreateEgressNotProtocolRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules that deny traffic to specified protocols for Calico. It processes the 'To' field in the intent request, which lists the protocols that should be denied. The function extracts protocol details from the 'To' field and creates rules to block egress traffic using these protocols, thus controlling egress traffic effectively by excluding certain protocols."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "protocol" && to.Protocol != "" {
			protocol := numorstring.ProtocolFromString(to.Protocol)
			calicoProtocol := convertProtocol(protocol)
			Rule := calico.Rule{
				NotProtocol: &calicoProtocol,
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// convertProtocol converts a numorstring.Protocol to a calicoapilib.Protocol
func convertProtocol(protocol numorstring.Protocol) calicoapilib.Protocol {
	if protocol.Type == numorstring.NumOrStringNum {
		return calicoapilib.Protocol{
			Type:   calicoapilib.NumOrStringNum,
			NumVal: protocol.NumVal,
		}
	}
	return calicoapilib.Protocol{
		Type:   calicoapilib.NumOrStringString,
		StrVal: protocol.StrVal,
	}
}

// ----------------------------
// Calico Policy: CIDR
// ----------------------------

// CreateIngressCIDRRule generates an ingress rule with CIDR blocks
func CreateIngressCIDRNets(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules with CIDR blocks for Calico. It processes the 'From' field in the intent request, which specifies the source CIDR blocks that need to be allowed access. The function iterates over the CIDR blocks listed in the 'From' field and creates ingress rules to permit traffic from these IP ranges. This is essential for controlling traffic that originates from (or terminates at) IP addresses in any of the given subnets."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "cidr" && len(from.Args) > 0 {
			Rule := calico.Rule{
				Source: calico.EntityRule{
					Nets: from.Args,
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateIngressCIDRRule generates an ingress rule with CIDR blocks
func CreateIngressNotCIDRNets(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules that deny traffic from specified CIDR blocks for Calico. It processes the 'From' field in the intent request, which lists the CIDR blocks to be denied. The function creates ingress rules that block traffic from these IP ranges, thus preventing unauthorized access from specified network segments."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "cidr" && len(from.Args) > 0 {
			Rule := calico.Rule{
				Source: calico.EntityRule{
					NotNets: from.Args,
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateEgressCIDRRule generates an egress rule with CIDR blocks
func CreateEgressCIDRNets(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules with CIDR blocks for Calico. It processes the 'To' field in the intent request, which specifies the destination CIDR blocks that need to be allowed access. The function iterates over the CIDR blocks listed in the 'To' field and creates egress rules to permit traffic to these IP ranges. This allows for fine-grained control over egress traffic that originates from (or terminates at) IP addresses in any of the given subnets."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "cidr" && len(to.Args) > 0 {
			Rule := calico.Rule{
				Destination: calico.EntityRule{
					Nets: to.Args,
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateEgressCIDRRule generates an egress rule with CIDR blocks
func CreateEgressNotCIDRNets(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules that deny traffic to specified CIDR blocks for Calico. It processes the 'To' field in the intent request, which lists the CIDR blocks to be denied. The function creates egress rules that block traffic to these IP ranges, helping to prevent the risk of data leakage to unauthorized network segments."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "cidr" && len(to.Args) > 0 {
			Rule := calico.Rule{
				Destination: calico.EntityRule{
					NotNets: to.Args,
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// ----------------------------
// Calico Policy: Ports
// ----------------------------

// CreateIngressSinglePortRule generates ingress rules with specific ports
func CreateIngressSinglePortRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules with specific ports for Calico. It processes the 'From' field in the intent request, which specifies the source ports that need to be allowed access. The function creates ingress rules that permit traffic from the specified ports, ensuring that the traffic can reach its destination. This is particularly useful for controlling traffic based on port numbers."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "port" {
			Rule := calico.Rule{
				Action: "Allow",
			}
			if from.Port != "" {
				port, err := numorstring.PortFromString(from.Port)
				if err != nil {
					return nil, err
				}
				calicoPort := convertPort(port)
				Rule.Destination.Ports = []calicoapilib.Port{calicoPort}
			}
			if from.Protocol != "" {
				protocol := numorstring.ProtocolFromString(from.Protocol)
				calicoProtocol := convertProtocol(protocol)
				Rule.Protocol = &calicoProtocol
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateEgressSinglePortRule generates egress rules with specific ports
func CreateEgressSinglePortRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules with specific ports for Calico. It processes the 'To' field in the intent request, which specifies the destination ports that need to be allowed access. The function creates egress rules that permit traffic to the specified ports, ensuring that the traffic can reach its destination. This is useful for managing egress traffic based on port numbers."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "port" {
			Rule := calico.Rule{
				Action: "Allow",
			}
			if to.Port != "" {
				port, err := numorstring.PortFromString(to.Port)
				if err != nil {
					return nil, err
				}
				calicoPort := convertPort(port)
				Rule.Destination.Ports = []calicoapilib.Port{calicoPort}
			}
			if to.Protocol != "" {
				protocol := numorstring.ProtocolFromString(to.Protocol)
				calicoProtocol := convertProtocol(protocol)
				Rule.Protocol = &calicoProtocol
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// convertPort converts a numorstring.Port to a calicoapilib.Port
func convertPort(port numorstring.Port) calicoapilib.Port {
	ruleDescription = "This function converts a numorstring.Port to a calicoapilib.Port."
	return calicoapilib.Port{
		MinPort: port.MinPort,
		MaxPort: port.MaxPort,
	}
}

// ----------------------------
// Calico Policy: HTTP
// ----------------------------

// CreateHTTPRules generates HTTP rules from ActionPoints
func CreateHTTPRules(actionPoints []v1.ActionPoint) ([]calico.HTTPMatch, error) {
	ruleDescription = "This function generates HTTP Match rules for Calico. It processes the action points, which specify the HTTP methods and paths that need to be allowed. The function creates HTTP match rules that permit traffic matching the specified HTTP methods and HTTP paths. This allows for precise control over HTTP traffic, ensuring that only specific methods and paths are permitted."
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
	ruleDescription = "This function converts path strings to HTTPPath objects for Calico. It processes the list of paths and creates HTTPPath objects that match the specified paths."
	var httpPaths []calico.HTTPPath
	for _, path := range paths {
		httpPaths = append(httpPaths, calico.HTTPPath{Exact: path})
	}
	return httpPaths
}

// ----------------------------
// Calico Policy: ICMP
// ----------------------------

// CreateIngressICMPRule generates an ingress rule with ICMP settings
func CreateIngressICMPRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates ingress rules with ICMP settings for Calico. It processes the 'From' field in the intent request, which specifies the source ICMP settings that need to be allowed access. The function creates ingress rules that permit traffic matching the specific type and code of ICMP traffic. The ICMP type and code are extracted from the 'From' field, and the rule is configured to allow traffic that matches these ICMP settings. This is useful for controlling ingress traffic based on specific ICMP messages, such as Echo Request (ping) or Destination Unreachable messages."
	var Rules []calico.Rule

	for _, from := range intentRequest.Rule.From {
		if from.Kind == "icmp" && len(from.Args) >= 2 {
			icmpType, err := strconv.Atoi(from.Args[0])
			if err != nil {
				return nil, err
			}
			icmpCode, err := strconv.Atoi(from.Args[1])
			if err != nil {
				return nil, err
			}

			Rule := calico.Rule{
				Protocol: &calicoapilib.Protocol{StrVal: "ICMP"},
				ICMP: &calico.ICMPFields{
					Type: &icmpType,
					Code: &icmpCode,
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// CreateEgressICMPRule generates an egress rule with ICMP settings
func CreateEgressICMPRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates egress rules with ICMP settings for Calico. It processes the 'To' field in the intent request, which specifies the destination ICMP settings that need to be allowed access. The function creates egress rules that permit traffic matching the specific type and code of ICMP traffic. The ICMP type and code are extracted from the 'To' field, and the rule is configured to allow traffic that matches these ICMP settings. This helps manage egress traffic by allowing only specific ICMP messages to be sent, such as Echo Reply or Time Exceeded messages."
	var Rules []calico.Rule

	for _, to := range intentRequest.Rule.To {
		if to.Kind == "icmp" && len(to.Args) >= 2 {
			icmpType, err := strconv.Atoi(to.Args[0])
			if err != nil {
				return nil, err
			}
			icmpCode, err := strconv.Atoi(to.Args[1])
			if err != nil {
				return nil, err
			}

			Rule := calico.Rule{
				Protocol: &calicoapilib.Protocol{StrVal: "ICMP"},
				ICMP: &calico.ICMPFields{
					Type: &icmpType,
					Code: &icmpCode,
				},
			}
			Rules = append(Rules, Rule)
		}
	}
	return Rules, nil
}

// ----------------------------
// Calico Policy: Action
// ----------------------------

// CreateActionRule generates a rule with specific actions for Calico
func CreateActionRule(intentRequest v1.IntentRequest) ([]calico.Rule, error) {
	ruleDescription = "This function generates rules with specific actions for Calico. It processes the 'Action' field in the intent request, which specifies the actions that need to be applied. The function creates rules that apply the specified actions. The 'Action' field can include actions such as 'Allow', 'Deny', 'Log', and others that dictate how the traffic should be handled. This allows for detailed control over the behavior of network policies, specifying exactly how traffic should be treated."
	var Rules []calico.Rule

	for _, action := range intentRequest.Rule.Action {
		Rule := calico.Rule{
			Action: calico.Action(action), // Convert rune to calico.Action type
		}
		Rules = append(Rules, Rule)
	}
	return Rules, nil
}
