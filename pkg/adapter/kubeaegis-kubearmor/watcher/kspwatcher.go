// kspwatcher.go
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

	kubearmorv1 "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1"
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

func kspInformer() cache.SharedIndexInformer {
	kspGvr := schema.GroupVersionResource{
		Group:    "security.kubearmor.com",
		Version:  "v1",
		Resource: "kubearmorpolicies",
	}
	informer := factory.ForResource(kspGvr).Informer()
	return informer
}

// WatchKsps watches for KSP events and logs deletions.
func WatchKsps(ctx context.Context, logger logr.Logger) {
	informer := kspInformer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			ksp, ok := obj.(*kubearmorv1.KubeArmorPolicy)
			if !ok {
				logger.Error(nil, "Could not cast to KubeArmorPolicy", "obj", obj)
				return
			}
			logger.Info("KubeArmorPolicy deleted", "KubeArmor.Name", ksp.Name, "KubeArmor.Namespace", ksp.Namespace)
		},
	})

	logger.Info("Starting KSP watcher")
	informer.Run(ctx.Done())
}
