package manager

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/adapter/k8s"
	"github.com/cclab-inu/KubeAegis/pkg/adapter/kubeaegis-cilium/converter"
	"github.com/cclab-inu/KubeAegis/pkg/adapter/kubeaegis-cilium/enforcer"
	watcher "github.com/cclab-inu/KubeAegis/pkg/adapter/watcher"
	"github.com/cclab-inu/KubeAegis/pkg/statusmanager"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
)

var (
	scheme    = runtime.NewScheme()
	k8sClient client.Client
)

func init() {
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(ciliumv2.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
	k8sClient = k8s.NewOrDie(scheme)
}

func Run(ctx context.Context, logger logr.Logger, KapName string, KapNamespace string) (string, error) {
	var kspname string
	kap, err := watcher.GetKubeAegisPolicy(ctx, k8sClient, KapName, KapNamespace)
	if err != nil {
		return "", err
	}
	logger.Info("KubeAegisPolicy fetched", "KubeAegis.Name", kap.Name, "KubeAegis.Namespace", kap.Namespace)

	ksp, err := converter.Converter(ctx, k8sClient, logger, kap)
	if err != nil {
		return "", err
	}

	if kspname, err = enforcer.Enforcer(ctx, k8sClient, logger, ksp, kap); err != nil {
		return "", err
	}

	if err := statusmanager.UpdateKapStatusAfterPolicy(ctx, k8sClient, ksp.Name, KapName, KapNamespace); err != nil {
		logger.Error(err, "failed to update KubeAegisPolicy status", "KubeAegis.Name", KapName, "KubeAegis.Namespace", KapNamespace)
		return "", err
	}

	return kspname, nil
}
