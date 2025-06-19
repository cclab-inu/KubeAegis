package converter

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	processor "github.com/cclab-inu/KubeAegis/pkg/adapter/processor"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	"github.com/cilium/cilium/pkg/policy/api"
)

func Converter(ctx context.Context, k8sClient client.Client, logger logr.Logger, kap *v1.KubeAegisPolicy) (*ciliumv2.CiliumNetworkPolicy, error) {
	logger.Info("CiliumNetworkPolicy started to transfer")

	ciliumNetworkPolicy := &ciliumv2.CiliumNetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateCNPName(kap.Name),
			Namespace: kap.Namespace,
		},
		Spec: &api.Rule{},
	}

	// Assuming that there might be multiple IntentRequests, we iterate through them.
	for _, intentRequest := range kap.Spec.IntentRequest {
		matchLabels, err := extractSelector(ctx, k8sClient, kap.Namespace, intentRequest.Selector)
		if err != nil {
			logger.Error(err, "failed to extract selector")
			return nil, err
		}
		if len(matchLabels) == 0 {
			continue
		}

		ciliumNetworkPolicy.Spec.EndpointSelector = api.NewESFromMatchRequirements(matchLabels, nil)
		if intentRequest.Rule.Action == "Block" && len(intentRequest.Rule.To) > 0 {
			ciliumNetworkPolicy.Spec.EgressDeny, _ = getEgressDeny(intentRequest)
		} else if intentRequest.Rule.Action == "Block" && len(intentRequest.Rule.From) > 0 {
			ciliumNetworkPolicy.Spec.IngressDeny, _ = getIngressDeny(intentRequest)
		} else if intentRequest.Rule.Action == "Allow" && len(intentRequest.Rule.To) > 0 {
			logger.Info("Egress Rules will be added at a later date ...")
			//ciliumNetworkPolicy.Spec.Egress = getEgress(intentRequest.Rule.To)
		} else if intentRequest.Rule.Action == "Allow" && len(intentRequest.Rule.From) > 0 {
			ciliumNetworkPolicy.Spec.Ingress, _ = getIngress(intentRequest)
		}
	}
	logger.Info("CiliumPolicy converted")
	return ciliumNetworkPolicy, nil
}

// extractSelector extracts match labels from a Selector.
func extractSelector(ctx context.Context, k8sClient client.Client, namespace string, selector v1.Selector) (map[string]string, error) {
	matchLabels := make(map[string]string) // Initialize map for match labels.

	// Process CEL expressions.
	if len(selector.CEL) > 0 {
		celExpressions := selector.CEL
		celMatchLabels, err := processor.ProcessCEL(ctx, k8sClient, namespace, celExpressions)
		if err != nil {
			return nil, fmt.Errorf("error processing CEL: %v", err)
		}
		for k, v := range celMatchLabels {
			matchLabels[k] = v
		}
	}

	// Process Match fields.
	if len(selector.Match) > 0 {
		processedMatchLabels, err := processor.ProcessMatchLabels(selector.Match)
		if err != nil {
			return nil, errors.Wrap(err, "error processing matchLabels")
		}
		for key, value := range processedMatchLabels {
			matchLabels[key] = value
		}
	}

	return matchLabels, nil
}

func generateCNPName(kapName string) string {
	return "cnp-" + kapName
}
