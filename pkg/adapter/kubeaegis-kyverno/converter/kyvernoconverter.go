package converter

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/adapter/processor"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
)

func Converter(ctx context.Context, k8sClient client.Client, logger logr.Logger, kap *v1.KubeAegisPolicy) (*kyvernov1.ClusterPolicy, error) {
	logger.Info("Converting KubeAegisPolicy to KyvernoPolicy")
	background := true
	hasValidateSubType := false
	kyvernoPolicy := &kyvernov1.ClusterPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeKypName(kap.Name),
			Namespace: kap.Namespace,
		},
		Spec: kyvernov1.Spec{
			Rules:                   []kyvernov1.Rule{},
			Background:              &background,
			ValidationFailureAction: convertActionToKyvernoFormat(kap.Spec.IntentRequest[0].Rule.Action),
		},
	}

	for _, intentRequest := range kap.Spec.IntentRequest {
		match, err := processor.ProcessMatch(intentRequest.Selector)
		if err != nil {
			logger.Error(err, "failed to extract match")
			return nil, err
		}
		if len(match.Any) == 0 && len(match.All) == 0 {
			continue
		}

		rule := kyvernov1.Rule{
			MatchResources: match,
			Name:           makeKypName(kap.Name) + "-" + intentRequest.Type,
		}
		var hasOperation bool
		for _, actionPoint := range intentRequest.Rule.ActionPoint {
			switch actionPoint.SubType {
			case "mutate":
				hasOperation = true
				mutation := handleMutate(ctx, intentRequest)
				if mutation != nil {
					rule.Mutation = mutation
				}
			case "validate":
				hasOperation = true
				validation := handleValidate(intentRequest)
				if validation != nil {
					rule.Validation = validation
				}
			case "verifyImage":
				hasOperation = true
				verifyImage := handleVerifyImage(intentRequest)
				if verifyImage != nil {
					rule.VerifyImages = append(rule.VerifyImages, *verifyImage)
				}
			}
		}
		if !hasOperation {
			logger.Error(nil, "No operation defined in the rule", "rule", rule.Name)
			continue
		}
		kyvernoPolicy.Spec.Rules = append(kyvernoPolicy.Spec.Rules, rule)
	}

	if hasValidateSubType {
		kyvernoPolicy.Spec.ValidationFailureAction = convertActionToKyvernoFormat(kap.Spec.IntentRequest[0].Rule.Action)
	}

	logger.Info("KyvernoPolicy converted")
	return kyvernoPolicy, nil
}

func makeKypName(kapName string) string {
	return "kyverno-" + kapName
}

func convertActionToKyvernoFormat(actionType string) kyvernov1.ValidationFailureAction {
	switch strings.ToLower(actionType) {
	case "enforce":
		return kyvernov1.ValidationFailureAction("Enforce")
	case "audit":
		return kyvernov1.ValidationFailureAction("Audit")
	default:
		return kyvernov1.ValidationFailureAction("Audit")
	}
}
