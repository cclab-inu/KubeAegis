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

	sample "importgopkg"
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

func sampleInformer() cache.SharedIndexInformer {
	PolicyGvr := schema.GroupVersionResource{
		Group:    SampleGroupString,
		Version:  SampleVersionString,
		Resource: SampleKindString,
	}
	informer := factory.ForResource(PolicyGvr).Informer()
	return informer
}

// Watchciliums watches for cilium events and logs deletions.
func WatchRealPolicy(ctx context.Context, logger logr.Logger) {
	informer := sampleInformer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			policy, ok := obj.(*sample.SampleSpecNamesKind)
			if !ok {
				logger.Error(nil, "Could not cast to SampleResourcePolicy", "obj", obj)
				return
			}
			logger.Info("SampleResourcePolicy deleted", "Policy.Name", policy.Name, "Policy.Namespace", policy.Namespace)
		},
	})

	logger.Info("Starting SampleResourcePolicy watcher")
	informer.Run(ctx.Done())
}
