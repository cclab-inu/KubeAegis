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

	calico "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
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

func calicoInformer() cache.SharedIndexInformer {
	PolicyGvr := schema.GroupVersionResource{
		Group:    "crd.projectcalico.org",
		Version:  "v1",
		Resource: "NetworkPolicy",
	}
	informer := factory.ForResource(PolicyGvr).Informer()
	return informer
}

// Watchciliums watches for cilium events and logs deletions.
func WatchRealPolicy(ctx context.Context, logger logr.Logger) {
	informer := calicoInformer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			policy, ok := obj.(*calico.NetworkPolicy)
			if !ok {
				logger.Error(nil, "Could not cast to NetworkPolicy", "obj", obj)
				return
			}
			logger.Info("NetworkPolicy deleted", "Policy.Name", policy.Name, "Policy.Namespace", policy.Namespace)
		},
	})

	logger.Info("Starting NetworkPolicy watcher")
	informer.Run(ctx.Done())
}
