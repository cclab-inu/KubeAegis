package watcher

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"

	tetragon "importgopkg"
)

var (
	factory dynamicinformer.DynamicSharedInformerFactory
)

func init() {
	k8sClient, err := dynamic.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		runtime.HandleError(err)
		return
	}
	factory = dynamicinformer.NewDynamicSharedInformerFactory(k8sClient, time.Minute)
}

func tetragonInformer() cache.SharedIndexInformer {
	PolicyGvr := schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v1alpha1",
		Resource: "TracingPolicyNamespaced",
	}
	informer := factory.ForResource(PolicyGvr).Informer()
	return informer
}

// Watchciliums watches for cilium events and logs deletions.
func WatchRealPolicy(ctx context.Context, logger logr.Logger) {
	informer := tetragonInformer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			policy, ok := obj.(*tetragon.TracingPolicyNamespaced)
			if !ok {
				logger.Error(nil, "Could not cast to TracingPolicyNamespaced", "obj", obj)
				return
			}
			logger.Info("TracingPolicyNamespaced deleted", "Policy.Name", policy.Name, "Policy.Namespace", policy.Namespace)
		},
	})

	logger.Info("Starting TracingPolicyNamespaced watcher")
	informer.Run(ctx.Done())
}
