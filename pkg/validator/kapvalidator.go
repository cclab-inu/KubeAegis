package validator

import (
	"context"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func KapValidator(ctx context.Context, k8sClient client.Client, logger logr.Logger, kap *v1.KubeAegisPolicy) ([]error, error) {
	var validationErrors []error
	logger.Info("Step 1: Check for the existence of a resource")
	if errs := ValidateExistence(ctx, k8sClient, kap); len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
		if len(validationErrors) > 0 {
			return validationErrors, nil
		}
	}

	logger.Info("Step 2: Check resource status and properties")
	if errs := ValidatePrecondition(ctx, k8sClient, kap); len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
		if len(validationErrors) > 0 {
			return validationErrors, nil
		}
	}

	return validationErrors, nil
}
