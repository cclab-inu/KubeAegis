// ciliumwatcher.go
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

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
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

func ciliumInformer() cache.SharedIndexInformer {
	ciliumGvr := schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v2",
		Resource: "CiliumNetworkPolicy",
	}
	informer := factory.ForResource(ciliumGvr).Informer()
	return informer
}

// Watchciliums watches for cilium events and logs deletions.
func Watchciliums(ctx context.Context, logger logr.Logger) {
	informer := ciliumInformer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			cilium, ok := obj.(*ciliumv2.CiliumNetworkPolicy)
			if !ok {
				logger.Error(nil, "Could not cast to CiliumPolicy", "obj", obj)
				return
			}
			logger.Info("CiliumPolicy deleted", "Cilium.Name", cilium.Name, "Cilium.Namespace", cilium.Namespace)
		},
	})

	logger.Info("Starting Cilium watcher")
	informer.Run(ctx.Done())
}
