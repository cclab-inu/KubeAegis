package processor

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
)

// ProcessMatchLabels processes any/all fields to generate matchLabels.
func ProcessMatchLabels(match []v1.Match) (map[string]string, error) {
	matchLabels := make(map[string]string)

	// Process logic for Any field.
	for _, match := range match {
		for key, value := range match.MatchLabels {
			matchLabels[key] = value
		}
	}

	return matchLabels, nil
}

func ProcessCEL(ctx context.Context, k8sClient client.Client, namespace string, expressions []string) (map[string]string, error) {
	logger := log.FromContext(ctx)

	// Retrieve pod list
	var podList corev1.PodList
	if err := k8sClient.List(ctx, &podList, client.InNamespace(namespace)); err != nil {
		logger.Error(err, "Error listing pods in namespace", "Namespace", namespace)
		return nil, fmt.Errorf("error listing pods: %v", err)
	}

	matchLabels := make(map[string]string)

	// Parse and evaluate label expressions
	for _, expr := range expressions {
		isNegated := checkNegation(expr)
		expr = preprocessExpression(expr)

		labels := extractLabelsFromExpression(expr, podList, isNegated)
		for k, v := range labels {
			matchLabels[k] = v
		}
	}
	return matchLabels, nil
}

func ProcessMatch(selector v1.Selector) (kyvernov1.MatchResources, error) {
	var match kyvernov1.MatchResources
	if len(selector.Match) > 0 {
		switch selector.Match[0].Condition {
		case "any":
			match.Any = make([]kyvernov1.ResourceFilter, 0)
		case "all":
			match.All = make([]kyvernov1.ResourceFilter, 0)
		default:
			return match, fmt.Errorf("unknown scope: %s", selector.Match[0].Condition)
		}
	}

	processedMatchLabels, _ := ProcessMatchLabels(selector.Match)

	for _, m := range selector.Match {
		resourceFilter := kyvernov1.ResourceFilter{
			ResourceDescription: kyvernov1.ResourceDescription{
				Kinds:      []string{m.Kind},
				Namespaces: []string{m.Namespace},
				Name:       m.Name,
				Selector: &metav1.LabelSelector{
					MatchLabels: processedMatchLabels,
				},
			},
		}

		if m.Condition == "any" {
			match.Any = append(match.Any, resourceFilter)
		} else if m.Condition == "all" {
			match.All = append(match.All, resourceFilter)
		}
	}
	return match, nil
}
