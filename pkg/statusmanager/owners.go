package statusmanager

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	karmorv1 "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
)

// SetOwnerReferences sets the KubeAegisPolicy as owner of the CiliumNetworkPolicy
func SetOwnerReferencesCNP(ctx context.Context, k8sClient client.Client, kap *v1.KubeAegisPolicy, cnp *ciliumv2.CiliumNetworkPolicy) error {
	// Fetch the KubeAegisPolicy to use it as an owner reference
	ownerRef := metav1.OwnerReference{
		APIVersion: "cclab.kubeaegis.com/v1",
		Kind:       "KubeAegisPolicy",
		Name:       kap.Name,
		UID:        kap.UID,
	}

	// Set the KubeAegisPolicy as the owner of the CiliumNetworkPolicy
	cnp.SetOwnerReferences(append(cnp.GetOwnerReferences(), ownerRef))
	return nil
}

// SetOwnerReferences sets the KubeAegisPolicy as owner of the KubeArmorPolicy
func SetOwnerReferencesKSP(ctx context.Context, k8sClient client.Client, kap *v1.KubeAegisPolicy, ksp *karmorv1.KubeArmorPolicy) error {
	// Fetch the KubeAegisPolicy to use it as an owner reference
	ownerRef := metav1.OwnerReference{
		APIVersion: "cclab.kubeaegis.com/v1",
		Kind:       "KubeAegisPolicy",
		Name:       kap.Name,
		UID:        kap.UID,
	}

	// Set the KubeAegisPolicy as the owner of the KubeArmorPolicy
	ksp.SetOwnerReferences(append(ksp.GetOwnerReferences(), ownerRef))

	return nil
}

// SetOwnerReferences sets the KubeAegisPolicy as owner of the KyvernoPolicy
func SetOwnerReferencesKCP(ctx context.Context, k8sClient client.Client, kap *v1.KubeAegisPolicy, kvp *kyvernov1.ClusterPolicy) error {
	// Fetch the KubeAegisPolicy to use it as an owner reference
	ownerRef := metav1.OwnerReference{
		APIVersion: "kyverno.io",
		Kind:       "ClusterPolicy",
		Name:       kap.Name,
		UID:        kap.UID,
	}

	// Set the KubeAegisPolicy as the owner of the KyvernoPolicy
	kvp.SetOwnerReferences(append(kvp.GetOwnerReferences(), ownerRef))
	return nil
}

// SetOwnerReferencesPod sets the AP as the owner of the Pod
func SetOwnerReferencesPodfromKSP(ctx context.Context, k8sClient client.Client, ksp *karmorv1.KubeArmorPolicy, pod *corev1.Pod) error {
	// Fetch the KubeAegisPolicy to use it as an owner reference
	ownerRef := metav1.OwnerReference{
		APIVersion: "security.kubearmor.com",
		Kind:       "kubearmorpolicies",
		Name:       ksp.Name,
		UID:        ksp.UID,
	}

	// Set the KubeAegisPolicy as the owner of the KubeArmorPolicy
	pod.SetOwnerReferences(append(ksp.GetOwnerReferences(), ownerRef))
	k8sClient.Update(ctx, pod)
	return nil
}

// SetOwnerReferencesPod sets the AP as the owner of the Pod
func SetOwnerReferencesPod(ctx context.Context, k8sClient client.Client, ap client.Object, pod *corev1.Pod) error {
	apKey := types.NamespacedName{Name: ap.GetName(), Namespace: ap.GetNamespace()}
	if err := k8sClient.Get(ctx, apKey, ap); err != nil {
		return err
	}

	ownerRef := metav1.OwnerReference{
		APIVersion: ap.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		Kind:       ap.GetObjectKind().GroupVersionKind().Kind,
		Name:       ap.GetName(),
		UID:        ap.GetUID(),
	}
	pod.SetOwnerReferences(append(pod.GetOwnerReferences(), ownerRef))
	return k8sClient.Update(ctx, pod)
}

// SetOwnerReferences sets the KubeAegisPolicy as owner of any resource
func SetOwnerReferences(ctx context.Context, k8sClient client.Client, kap *v1.KubeAegisPolicy, resource client.Object, apiVersion string, kind string) error {
	// Fetch the KubeAegisPolicy to use it as an owner reference
	ownerRef := metav1.OwnerReference{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       kap.Name,
		UID:        kap.UID,
	}

	// Set the KubeAegisPolicy as the owner of the resource
	resource.SetOwnerReferences(append(resource.GetOwnerReferences(), ownerRef))
	return k8sClient.Update(ctx, resource)
}
