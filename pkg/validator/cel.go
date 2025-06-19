package validator

import (
	"context"
	"fmt"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ValidateCEL processes CEL expressions and validates if the resources meeting the conditions actually exist.
func ValidateCEL(ctx context.Context, k8sClient client.Client, kap *v1.KubeAegisPolicy) (bool, error) {
	for _, intentRequest := range kap.Spec.IntentRequest {
		if len(intentRequest.Selector.CEL) > 0 {
			env, err := cel.NewEnv(
				cel.Declarations(
					decls.NewVar("labels", decls.NewMapType(decls.String, decls.String)),
				),
			)
			if err != nil {
				return false, errors.Wrap(err, "failed to create CEL environment")
			}

			for _, expr := range intentRequest.Selector.CEL {
				ast, issues := env.Compile(expr)
				if issues != nil && issues.Err() != nil {
					return false, fmt.Errorf("CEL compile error: %s", issues.Err())
				}

				prg, err := env.Program(ast)
				if err != nil {
					return false, fmt.Errorf("failed to create CEL program: %w", err)
				}

				var podList corev1.PodList
				if err := k8sClient.List(ctx, &podList, client.InNamespace(kap.Namespace)); err != nil {
					return false, fmt.Errorf("error listing pods: %w", err)
				}

				resourceFound := false
				for _, pod := range podList.Items {
					out, _, err := prg.Eval(map[string]interface{}{
						"labels": pod.GetLabels(),
					})
					if err != nil {
						return false, fmt.Errorf("CEL evaluation error: %w", err)
					}

					if out.Value().(bool) {
						resourceFound = true
						break
					}
				}

				if !resourceFound {
					return false, fmt.Errorf("no resources found matching the CEL expression: %s", expr)
				}
			}
		}
	}

	return true, nil
}
