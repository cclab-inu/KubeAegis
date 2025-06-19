package statusmanager

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/reporter"
)

const (
	StatusCreated = "Created"
)

func UpdateKapStatus(ctx context.Context, k8sClient client.Client, kapName, kapNamespace string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		kap := &v1.KubeAegisPolicy{}
		if err := k8sClient.Get(ctx, types.NamespacedName{Name: kapName, Namespace: kapNamespace}, kap); err != nil {
			return err
		}

		kap.Status.LastUpdated = metav1.Now()
		kap.Status.Status = StatusCreated

		return k8sClient.Status().Update(ctx, kap)
	})
}

func UpdateKapStatusAfterPolicy(ctx context.Context, k8sClient client.Client, currPolicyFullName, kapName, namespace string) error {
	if retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		latestKap := &v1.KubeAegisPolicy{}
		if err := k8sClient.Get(ctx, types.NamespacedName{Name: kapName, Namespace: namespace}, latestKap); err != nil {
			return nil
		}

		updateandcountAPinfo(latestKap, currPolicyFullName)
		if err := k8sClient.Status().Update(ctx, latestKap); err != nil {
			return err
		}

		return nil
	}); retryErr != nil {
		return retryErr
	}

	return nil
}

// UpdateKapStatusAfterPolicy updates the provided KubeAegisPolicy status with the number and
// names of its descendant policies that were created.
func UpdateKapStatusAfterPolicywithResource(ctx context.Context, k8sClient client.Client, currPolicyFullName, kapName, namespace string, resourceNames []string) error {
	// Since multiple adapters may attempt to update the KubeAegisPolicy status
	// concurrently, potentially leading to conflicts. To ensure data consistency,
	// retry on write failures. On conflict, the update is retried with an
	// exponential backoff strategy. This provides resilience against potential
	// issues while preventing indefinite retries in case of persistent conflicts.
	if retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		latestKap := &v1.KubeAegisPolicy{}
		if err := k8sClient.Get(ctx, types.NamespacedName{Name: kapName, Namespace: namespace}, latestKap); err != nil {
			return err
		}

		updateandcountAPinfo(latestKap, currPolicyFullName)
		for _, resourceName := range resourceNames {
			kindResourceName := fmt.Sprintf("Pod/%s", resourceName)
			updateandcountResourceinfo(latestKap, kindResourceName)
		}
		if err := k8sClient.Status().Update(ctx, latestKap); err != nil {
			return err
		}

		return nil
	}); retryErr != nil {
		return retryErr
	}

	return nil
}

func updateandcountAPinfo(latestKap *v1.KubeAegisPolicy, currPolicyFullName string) {
	if !contains(latestKap.Status.ListofAPs, currPolicyFullName) {
		latestKap.Status.NumberOfAPs++
		latestKap.Status.ListofAPs = append(latestKap.Status.ListofAPs, currPolicyFullName)
	}
}

func updateandcountResourceinfo(latestKap *v1.KubeAegisPolicy, currResourceFullName string) {
	if !contains(latestKap.Status.ListofResources, currResourceFullName) {
		latestKap.Status.NumberOfResources++
		latestKap.Status.ListofResources = append(latestKap.Status.ListofResources, currResourceFullName)
	}
}

func contains(existingPolicies []string, policy string) bool {
	for _, existingPolicy := range existingPolicies {
		if existingPolicy == policy {
			return true
		}
	}
	return false
}

func NotifyReporter(ctx context.Context, kap *v1.KubeAegisPolicy, configMap corev1.ConfigMap, adapterPolicy string) error {
	logger := log.FromContext(ctx)
	if kap.Spec.EnableReporting {
		if err := reporter.GenerateReport(ctx, logger, kap, configMap, adapterPolicy); err != nil {
			logger.Error(err, "failed to generate report")
			return err
		}
	}
	return nil
}
