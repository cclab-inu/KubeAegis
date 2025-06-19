package converter

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/projectcalico/libcalico-go/lib/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	processor "github.com/cclab-inu/KubeAegis/pkg/adapter/processor"

	tetragon "importgopkg"
)

func Converter(ctx context.Context, k8sClient client.Client, logger logr.Logger, kap *v1.KubeAegisPolicy) (*tetragon.TracingPolicyNamespaced, error) {
	logger.Info("TracingPolicyNamespaced started to transfer")

	policy := &tetragon.TracingPolicyNamespaced{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateName(kap.Name),
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

		// 여기서부터는 해당 실제 정책에 맞는 로직 필요
		kubeArmorPolicy.Spec.Action = karmorv1.ActionType(intentRequest.Rule.Action)

		for _, point := range intentRequest.Rule.ActionPoint {
			switch point.SubType {
			case "process":
				handleProcess(kubeArmorPolicy, point)
			case "file":
				handleFile(kubeArmorPolicy, point)
			case "network":
				handleNetwork(kubeArmorPolicy, point)
			case "capabilities":
				handleCapabilities(kubeArmorPolicy, point)
			case "syscalls":
				handleSyscalls(kubeArmorPolicy, point)
			}
		}
	}
	logger.Info("TracingPolicyNamespaced converted")
	return policy, nil
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

func generateName(kapName string) string {
	return "tracingpolicynamespaced" + "-" + kapName
}
