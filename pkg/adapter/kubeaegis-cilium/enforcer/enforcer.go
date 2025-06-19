package enforcer

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/statusmanager"
	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
)

func Enforcer(ctx context.Context, k8sClient client.Client, logger logr.Logger, cnp *ciliumv2.CiliumNetworkPolicy, kap *v1.KubeAegisPolicy) (string, error) {
	// Check if the policy already exists
	existingPolicy := &ciliumv2.CiliumNetworkPolicy{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: cnp.Name, Namespace: cnp.Namespace}, existingPolicy)
	if err != nil && !apierrors.IsNotFound(err) {
		logger.Error(err, "failed to fetch CiliumNetworkPolicy", "Cilium.Name", cnp.Name, "Cilium.Namespace", cnp.Namespace)
		return "", err
	}

	if err := statusmanager.SetOwnerReferencesCNP(ctx, k8sClient, kap, cnp); err != nil {
		logger.Error(err, "failed to set KubeAegisPolicy as owner of CiliumNetworkPolicy")
		return "", err
	}

	// Update if exists, create otherwise
	if apierrors.IsNotFound(err) {
		logger.Info("CiliumNetworkPolicy enforced", "Cilium.Name", cnp.Name, "Cilium.Namespace", cnp.Namespace)
		if err := k8sClient.Create(ctx, cnp); err != nil {
			logger.Error(err, "failed to create CiliumNetworkPolicy", "Cilium.Name", cnp.Name, "Cilium.Namespace", cnp.Namespace)
			return "", err
		}
	} else {
		logger.Info("CiliumNetworkPolicy updated", "PolicyName", cnp.Name, "Cilium.Namespace", cnp.Namespace)
		existingPolicy.Spec = cnp.Spec
		if err := k8sClient.Update(ctx, existingPolicy); err != nil {
			logger.Error(err, "failed to update CiliumNetworkPolicy", "Cilium.Name", cnp.Name, "Cilium.Namespace", cnp.Namespace)
			return "", err
		}
	}

	return cnp.Name, nil
}
