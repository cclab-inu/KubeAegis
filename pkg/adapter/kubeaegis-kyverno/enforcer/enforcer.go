package enforcer

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/statusmanager"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
)

func Enforcer(ctx context.Context, k8sClient client.Client, logger logr.Logger, kyvernoPolicy *kyvernov1.ClusterPolicy, kap *v1.KubeAegisPolicy) (string, error) {
	// Check if the policy already exists
	existingPolicy := &kyvernov1.ClusterPolicy{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: kyvernoPolicy.Name, Namespace: kyvernoPolicy.Namespace}, existingPolicy)
	if err != nil && !apierrors.IsNotFound(err) {
		logger.Error(err, "failed to fetch KyvernoPolicy", "Kyverno.Name", kyvernoPolicy.Name, "Kyverno.Namespace", kyvernoPolicy.Namespace)
		return "", err
	}

	if err := statusmanager.SetOwnerReferencesKCP(ctx, k8sClient, kap, kyvernoPolicy); err != nil {
		logger.Error(err, "failed to set KubeAegisPolicy as owner of KyvernoPolicy")
		return "", err
	}

	// Update if exists, create otherwise
	if apierrors.IsNotFound(err) {
		logger.Info("KyvernoPolicy enforced", "Kyverno.Name", kyvernoPolicy.Name, "Kyverno.Namespace", kyvernoPolicy.Namespace)
		if err := k8sClient.Create(ctx, kyvernoPolicy); err != nil {
			logger.Error(err, "failed to create KyvernoPolicy", "Kyverno.Name", kyvernoPolicy.Name, "Kyverno.Namespace", kyvernoPolicy.Namespace)
			return "", err
		}
	} else {
		logger.Info("KyvernoPolicy updated", "PolicyName", kyvernoPolicy.Name, "Kyverno.Namespace", kyvernoPolicy.Namespace)
		existingPolicy.Spec = kyvernoPolicy.Spec
		if err := k8sClient.Update(ctx, existingPolicy); err != nil {
			logger.Error(err, "failed to update KyvernoPolicy", "Kyverno.Name", kyvernoPolicy.Name, "Kyverno.Namespace", kyvernoPolicy.Namespace)
			return "", err
		}
	}

	return kyvernoPolicy.Name, nil
}
