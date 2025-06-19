// Existence Validator
// Validate the existence of a Kubernetes resource for a given selector
// and validate the state of the namespace to ensure that the resource exists and its state is valid.
package validator

import (
	"context"
	"fmt"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	processor "github.com/cclab-inu/KubeAegis/pkg/adapter/processor"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ValidateExistence(ctx context.Context, k8sClient client.Client, kap *v1.KubeAegisPolicy) []error {
	var errs []error

	// Iterate over each intentRequest to check resource existence
	for _, intentRequest := range kap.Spec.IntentRequest {
		if (intentRequest.Type == "system" && len(intentRequest.Rule.From) > 0) || (intentRequest.Type == "system" && len(intentRequest.Rule.To) > 0) {
			errs = append(errs, errors.Errorf("Most likely you want to create a policy of network type, check the 'type' of the policy you are creating.: %v", intentRequest.Type))
		}

		matchLabels, err := processor.ProcessMatchLabels(intentRequest.Selector.Match)
		if err != nil {
			errs = append(errs, fmt.Errorf("error processing match labels: %v", err))
			continue
		}
		selector := labels.SelectorFromSet(matchLabels)

		if len(intentRequest.Selector.CEL) > 0 {
			if _, err := ValidateCEL(ctx, k8sClient, kap); err != nil {
				errs = append(errs, err)
			}
		}

		if len(intentRequest.Selector.Match) > 0 {
			for _, match := range intentRequest.Selector.Match {

				if err := validateNamespaces(ctx, k8sClient, match); err != nil {
					errs = append(errs, err)
				}
				err := validateNamespaceStatus(ctx, k8sClient, match.Namespace)
				if err != nil {
					errs = append(errs, err)
				}

				switch match.Kind {
				case "Pod":
					err := checkPodExistence(ctx, k8sClient, intentRequest.Selector.Match[0].Namespace, selector)
					if err != nil {
						errs = append(errs, err)
					}
				case "Service":
					if err := validateServices(ctx, k8sClient, match); err != nil {
						errs = append(errs, err)
					}
				case "Deployment":
					err := validateDeployments(ctx, k8sClient, match)
					if err != nil {
						errs = append(errs, err)
					}
				case "ConfigMap":
					err := validateConfigMaps(ctx, k8sClient, match.Namespace, selector)
					if err != nil {
						errs = append(errs, err)
					}
				}
			}
		}
	}

	return errs
}

// validatePods checks if the Pods matching the selector exist.
func checkPodExistence(ctx context.Context, k8sClient client.Client, namespace string, selector labels.Selector) error {
	var pods corev1.PodList
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabelsSelector{Selector: selector},
	}

	if err := k8sClient.List(ctx, &pods, listOpts...); err != nil {
		return fmt.Errorf("error listing pods: %v", err)
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no matching pods found in namespace %s with labels %v", namespace, selector)
	}

	return nil
}

// validateServices checks if the Services matching the selector exist.
func validateServices(ctx context.Context, k8sClient client.Client, match v1.Match) error {
	var services corev1.ServiceList
	listOpts := []client.ListOption{
		client.InNamespace(match.Namespace),
		client.MatchingLabels(match.MatchLabels),
	}
	if err := k8sClient.List(ctx, &services, listOpts...); err != nil {
		return errors.Errorf("error fetching services: %v", err)
	}
	if len(services.Items) == 0 {
		return errors.Errorf("no services found matching the selector")
	}
	return nil
}

func validateDeployments(ctx context.Context, k8sClient client.Client, match v1.Match) error {
	var deployments appsv1.DeploymentList
	listOpts := []client.ListOption{
		client.InNamespace(match.Namespace),
		client.MatchingLabels(match.MatchLabels),
	}
	if err := k8sClient.List(ctx, &deployments, listOpts...); err != nil {
		return errors.Errorf("error fetching deployments: %v", err)
	}
	if len(deployments.Items) == 0 {
		return errors.Errorf("no deployments found matching the selector")
	}
	return nil
}

// validateConfigMaps checks if the ConfigMaps matching the selector exist.
func validateConfigMaps(ctx context.Context, k8sClient client.Client, namespace string, selector labels.Selector) error {
	var configMaps corev1.ConfigMapList
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabelsSelector{Selector: selector},
	}

	if err := k8sClient.List(ctx, &configMaps, listOpts...); err != nil {
		return fmt.Errorf("error listing ConfigMaps: %v", err)
	}

	if len(configMaps.Items) == 0 {
		return fmt.Errorf("no matching ConfigMaps found in namespace %s with labels %v", namespace, selector)
	}

	return nil
}

func validateNamespaces(ctx context.Context, k8sClient client.Client, match v1.Match) error {
	var namespace corev1.Namespace
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: match.Namespace}, &namespace); err != nil {
		return errors.Errorf("error fetching namespace '%s': %v", match.Namespace, err)
	}
	return nil
}

// validateNamespaceStatus checks the status of the namespace to ensure it is active.
func validateNamespaceStatus(ctx context.Context, k8sClient client.Client, namespace string) error {
	var ns corev1.Namespace
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: namespace}, &ns); err != nil {
		return errors.Errorf("error fetching namespace '%s': %v", namespace, err)
	}

	if ns.Status.Phase != corev1.NamespaceActive {
		return errors.Errorf("namespace '%s' is not in an active phase", namespace)
	}

	return nil
}
