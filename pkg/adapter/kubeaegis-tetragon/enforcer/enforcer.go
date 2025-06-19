package enforcer

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tetragon "importgopkg"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/statusmanager"
)

func Enforcer(ctx context.Context, k8sClient client.Client, logger logr.Logger, policy *tetragon.TracingPolicyNamespaced, kap *v1.KubeAegisPolicy) (string, error) {
	// Check if the policy already exists
	existingPolicy := &tetragon.TracingPolicyNamespaced{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: policy.Name, Namespace: policy.Namespace}, existingPolicy)
	if err != nil && !apierrors.IsNotFound(err) {
		logger.Error(err, "failed to fetch TracingPolicyNamespaced", "Policy.Name", policy.Name, "Policy.Namespace", policy.Namespace)
		return "", err
	}

	apiversion := "cilium.io" + "/" + "v1alpha1"
	kind := "TracingPolicyNamespaced"
	if err := statusmanager.SetOwnerReferences(ctx, k8sClient, kap, policy, apiversion, kind); err != nil {
		logger.Error(err, "failed to set KubeAegisPolicy as owner of TracingPolicyNamespaced")
		return "", err
	}

	// Update if exists, create otherwise
	if apierrors.IsNotFound(err) {
		logger.Info("TracingPolicyNamespaced enforced", "Policy.Name", policy.Name, "Policy.Namespace", policy.Namespace)
		if err := k8sClient.Create(ctx, policy); err != nil {
			logger.Error(err, "failed to create TracingPolicyNamespaced", "Policy.Name", policy.Name, "Policy.Namespace", policy.Namespace)
			return "", err
		}
	} else {
		logger.Info("TracingPolicyNamespaced updated", "Policy.Name", policy.Name, "Policy.Namespace", policy.Namespace)
		existingPolicy.Spec = policy.Spec
		if err := k8sClient.Update(ctx, existingPolicy); err != nil {
			logger.Error(err, "failed to update TracingPolicyNamespaced", "Policy.Name", policy.Name, "Policy.Namespace", policy.Namespace)
			return "", err
		}
	}

	return policy.Name, nil
}
