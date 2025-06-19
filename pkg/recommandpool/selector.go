package recommandpool

import (
	"context"
	"fmt"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	processor "github.com/cclab-inu/KubeAegis/pkg/adapter/processor"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ----------------------------
// Selector (2)
// ----------------------------

// ExtractSelector extracts match labels from a Selector.
func ExtractSelector(ctx context.Context, k8sClient client.Client, namespace string, selector v1.Selector) (map[string]string, error) {
	ruleDescription = "This function extracts match labels from a Selector in a Kubernetes environment. It processes CEL expressions and match fields from the selector, generating a map of labels that match the specified criteria."
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

// ExtractStringSelector generates a selector string from given labels.
func ExtractGiveFormatsSelector(labels map[string]string) string {
	ruleDescription = "This function generates a selector string from given labels. It constructs a string representation of label-based selectors, enabling selection of resources that match the specified labels."
	var selector string
	for key, value := range labels {
		if selector != "" {
			selector += " && "
		}
		selector += fmt.Sprintf("%s == '%s'", key, value)
	}
	return selector
}
