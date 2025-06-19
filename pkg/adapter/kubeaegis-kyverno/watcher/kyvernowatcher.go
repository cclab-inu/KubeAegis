// kyvernowatcher.go
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

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
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

func kyvernoInformer() cache.SharedIndexInformer {
	kyvernoGvr := schema.GroupVersionResource{
		Group:    "kyverno.io",
		Version:  "v1",
		Resource: "ClusterPolicy",
	}
	informer := factory.ForResource(kyvernoGvr).Informer()
	return informer
}

/*func kyvernoCleanupInformer() cache.SharedIndexInformer {
	kyvernoGvr := schema.GroupVersionResource{
		Group:    "kyverno.io",
		Version:  "v2beta1",
		Resource: "ClusterCleanupPolicy",
	}
	informer := factory.ForResource(kyvernoGvr).Informer()
	return informer
}*/

// Watchkyvernos watches for kyverno events and logs deletions.
func WatchKyvernos(ctx context.Context, logger logr.Logger) {
	informer := kyvernoInformer()
	//cleanupInformoer := kyvernoCleanupInformer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			kyverno, ok := obj.(*kyvernov1.ClusterPolicy)
			if !ok {
				logger.Error(nil, "Could not cast to KyvernoPolicy", "obj", obj)
				return
			}
			logger.Info("KyvernoPolicy deleted", "Kyverno.Name", kyverno.Name, "Kyverno.Namespace", kyverno.Namespace)
		},
	})

	logger.Info("Starting Kyverno watcher")
	informer.Run(ctx.Done())
	//cleanupInformoer.Run(ctx.Done())
}
