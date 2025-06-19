package enforcer

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/statusmanager"
	karmorv1 "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1"
)

func Enforcer(ctx context.Context, k8sClient client.Client, logger logr.Logger, kubeArmorPolicy *karmorv1.KubeArmorPolicy, kap *v1.KubeAegisPolicy) (string, error) {
	// Check if the policy already exists
	existingPolicy := &karmorv1.KubeArmorPolicy{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: kubeArmorPolicy.Name, Namespace: kubeArmorPolicy.Namespace}, existingPolicy)
	if err != nil && !apierrors.IsNotFound(err) {
		logger.Error(err, "failed to fetch KubeArmorPolicy", "KubeArmor.Name", kubeArmorPolicy.Name, "KubeArmor.Namespace", kubeArmorPolicy.Namespace)
		return "", err
	}

	if err := statusmanager.SetOwnerReferencesKSP(ctx, k8sClient, kap, kubeArmorPolicy); err != nil {
		logger.Error(err, "failed to set KubeAegisPolicy as owner of KubeArmorPolicy")
		return "", err
	}

	// Update if exists, create otherwise
	if apierrors.IsNotFound(err) {
		logger.Info("KubeArmorPolicy enforced", "KubeArmor.Name", kubeArmorPolicy.Name, "KubeArmor.Namespace", kubeArmorPolicy.Namespace)
		if err := k8sClient.Create(ctx, kubeArmorPolicy); err != nil {
			logger.Error(err, "failed to create KubeArmorPolicy", "KubeArmor.Name", kubeArmorPolicy.Name, "KubeArmor.Namespace", kubeArmorPolicy.Namespace)
			return "", err
		}
	} else {
		logger.Info("KubeArmorPolicy updated", "PolicyName", kubeArmorPolicy.Name, "KubeArmor.Namespace", kubeArmorPolicy.Namespace)
		existingPolicy.Spec = kubeArmorPolicy.Spec
		if err := k8sClient.Update(ctx, existingPolicy); err != nil {
			logger.Error(err, "failed to update KubeArmorPolicy", "KubeArmor.Name", kubeArmorPolicy.Name, "KubeArmor.Namespace", kubeArmorPolicy.Namespace)
			return "", err
		}
	}

	return kubeArmorPolicy.Name, nil
}
