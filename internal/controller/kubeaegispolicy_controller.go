/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/exporter"
	"github.com/cclab-inu/KubeAegis/pkg/statusmanager"
	"github.com/cclab-inu/KubeAegis/pkg/validator"
)

// KubeAegisPolicyReconciler reconciles a KubeAegisPolicy object
type KubeAegisPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cclab.kubeaegis.com,resources=kubeaegispolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cclab.kubeaegis.com,resources=kubeaegispolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cclab.kubeaegis.com,resources=kubeaegispolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KubeAegisPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *KubeAegisPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	kap := &v1.KubeAegisPolicy{}
	err := r.Get(ctx, req.NamespacedName, kap)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("KubeAegis not found. Ignoring since object must be deleted")
			var configMap corev1.ConfigMap
			if err := r.Get(ctx, client.ObjectKey{Name: "adapter-config", Namespace: "default"}, &configMap); err == nil {
				if err := exporter.NotifyAdapterOfPolicyDeletion(ctx, req.NamespacedName, kap, configMap); err != nil {
					logger.Error(err, "Failed to notify adapter of KAP deletion")
				}
			}
			return doNotRequeue()
		}
		logger.Error(err, "failed to fetch KubeAegis", "KubeAegis.Name", req.Name)
		return requeueWithError(err)
	}
	logger.Info("KubeAegis found", "KubeAegis.Name", kap.Name, "KubeAegis.Namespace", kap.Namespace)

	validationErrors, err := validator.KapValidator(ctx, r.Client, logger, kap)
	if err != nil {
		return requeueWithError(err)
	}
	if len(validationErrors) > 0 {
		logger.Info("Identified a misconfiguration of KubeAegis", "ValidationErrors", validationErrors)
		return doNotRequeue()
	}
	logger.Info("KubeAegis verified", "validationErrors count", len(validationErrors))

	var configMap corev1.ConfigMap
	if err := r.Get(ctx, client.ObjectKey{Name: "adapter-config", Namespace: "default"}, &configMap); err != nil {
		logger.Error(err, "error fetching adapter config", "KubeAegis.Name", req.Name, "KubeAegis.Namespace", req.Namespace)
		return requeueWithError(err)
	}
	if err := exporter.DispatchPolicyToAdapters(ctx, r.Client, logger, kap, configMap); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("No adapters are currently found")
			return doNotRequeue()
		}
		logger.Info("failed to dispatch policy to adapters", "error", err)
		return doNotRequeue()
	}

	if err := statusmanager.UpdateKapStatus(ctx, r.Client, kap.Name, kap.Namespace); err != nil {
		logger.Error(err, "failed to update KubeAegisPolicy status")
		return requeueWithError(err)
	}

	return doNotRequeue()

}

// SetupWithManager sets up the controller with the Manager.
func (r *KubeAegisPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.KubeAegisPolicy{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)

}
