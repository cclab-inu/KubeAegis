package watcher

import (
	"context"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/pkg/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetKubeAegisPolicy(ctx context.Context, k8sClient client.Client, kapName string, kapNamespace string) (*v1.KubeAegisPolicy, error) {
	kap := &v1.KubeAegisPolicy{}
	key := client.ObjectKey{
		Namespace: kapNamespace,
		Name:      kapName,
	}

	if err := k8sClient.Get(ctx, key, kap); err != nil {
		return nil, errors.Wrapf(err, "failed fetch KubeAegisPolicy: %s", kapName)
	}
	return kap, nil
}
