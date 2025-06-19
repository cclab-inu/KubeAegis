package converter

import (
	"context"
	"encoding/json"
	"strings"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	admissionapi "k8s.io/pod-security-admission/api"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func handleMutate(ctx context.Context, intentRequest v1.IntentRequest) *kyvernov1.Mutation {
	logger := log.FromContext(ctx)
	mutation := &kyvernov1.Mutation{}
	for _, event := range intentRequest.Rule.ActionPoint {
		if event.SubType == "mutate" {
			if event.Resource.Kind == "annotations" {
				annotations := make(map[string]string)
				for _, detailMap := range event.Resource.Details {
					for key, value := range detailMap {
						annotations[key] = value
					}
				}
				annotationsPatch := map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": annotations,
					},
				}
				patchBytes, err := json.Marshal(annotationsPatch)
				if err != nil {
					logger.Error(err, "Failed to marshal annotations patch")
					continue
				}
				patch := apiextv1.JSON{Raw: patchBytes}
				mutation.RawPatchStrategicMerge = &patch
				return mutation
			} else if event.Resource.Kind == "label" {
				labels := make(map[string]string)
				for _, detailMap := range event.Resource.Details {
					for labelKey, labelValue := range detailMap {
						labels[labelKey] = labelValue
					}
				}
				labelsPatch := map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": labels,
					},
				}
				patchBytes, err := json.Marshal(labelsPatch)
				if err != nil {
					logger.Error(err, "Failed to marshal labels patch")
					continue
				}
				patch := apiextv1.JSON{Raw: patchBytes}
				mutation.RawPatchStrategicMerge = &patch
				return mutation
			}
		}
	}
	return nil
}

func handleValidate(intentRequest v1.IntentRequest) *kyvernov1.Validation {
	var validation *kyvernov1.Validation

	for _, point := range intentRequest.Rule.ActionPoint {
		if validation == nil {
			validation = &kyvernov1.Validation{}
		}

		switch point.SubType {
		case "cel":
			for _, detailMap := range point.Resource.Details {
				expression, exprOk := detailMap["expression"]
				message, msgOk := detailMap["message"]
				if exprOk && msgOk {
					validation.CEL = &kyvernov1.CEL{
						Expressions: []admissionv1.Validation{
							{
								Expression: expression,
								Message:    message,
							},
						},
					}
					return validation
				}
			}

		case "podSecurity":
			for _, detailMap := range point.Resource.Details {
				level, levelOk := detailMap["level"]
				version, versionOk := detailMap["version"]
				if levelOk && versionOk {
					var podSecurityLevel admissionapi.Level
					switch level {
					case "privileged":
						podSecurityLevel = admissionapi.LevelPrivileged
					case "baseline":
						podSecurityLevel = admissionapi.LevelBaseline
					case "restricted":
						podSecurityLevel = admissionapi.LevelRestricted
					default:
						continue
					}

					validation.PodSecurity = &kyvernov1.PodSecurity{
						Level:   podSecurityLevel,
						Version: version,
					}
					return validation
				}
			}

		case "deny":
			var tempConditions struct {
				Any []kyvernov1.Condition `json:"any,omitempty"`
				All []kyvernov1.Condition `json:"all,omitempty"`
			}

			for _, filter := range point.Resource.Filter {
				keyJSON, err := json.Marshal(filter.Key)
				if err != nil {
					continue
				}
				valueJSON, err := json.Marshal(filter.Value)
				if err != nil {
					continue
				}

				condition := kyvernov1.Condition{
					RawKey:   &apiextv1.JSON{Raw: keyJSON},
					Operator: kyvernov1.ConditionOperator(filter.Operator),
					RawValue: &apiextv1.JSON{Raw: valueJSON},
				}

				switch filter.Condition {
				case "any":
					tempConditions.Any = append(tempConditions.Any, condition)
				case "all":
					tempConditions.All = append(tempConditions.All, condition)
				}
			}
			conditionsJSON, err := json.Marshal(tempConditions)
			if err != nil {
				return nil
			}

			var cw kyvernov1.ConditionsWrapper
			if err := json.Unmarshal(conditionsJSON, &cw.Conditions); err != nil {
				return nil
			}

			validation.Deny = &kyvernov1.Deny{
				RawAnyAllConditions: &cw,
			}

		case "pattern":
			if len(point.Resource.Details) > 0 {
				patternMap := make(map[string]interface{})
				for _, detailMap := range point.Resource.Details {
					for k, v := range detailMap {
						patternMap[k] = v
					}
				}
				patternJSON, err := json.Marshal(patternMap)
				if err != nil {
					continue
				}
				validation.RawPattern = &apiextv1.JSON{Raw: patternJSON}
				return validation
			}
		}
	}
	return validation
}

func handleVerifyImage(intentRequest v1.IntentRequest) *kyvernov1.ImageVerification {
	imageVerification := &kyvernov1.ImageVerification{}

	for _, point := range intentRequest.Rule.ActionPoint {
		var details []string
		for _, detailMap := range point.Resource.Details {
			for key := range detailMap {
				details = append(details, key)
			}
		}
		imageVerification.ImageReferences = append(imageVerification.ImageReferences, details...)

		var attestorSet kyvernov1.AttestorSet
		if len(point.Resource.Keys) > 0 {
			for _, key := range point.Resource.Keys {
				if strings.HasPrefix(key, "kms:") {
					attestorSet.Entries = append(attestorSet.Entries, kyvernov1.Attestor{
						Keys: &kyvernov1.StaticKeyAttestor{KMS: key},
					})
				} else if strings.HasPrefix(key, "{{") {
					attestorSet.Entries = append(attestorSet.Entries, kyvernov1.Attestor{
						Keys: &kyvernov1.StaticKeyAttestor{PublicKeys: key},
					})
				}
			}
		}

		if len(point.Resource.Keyless) > 0 {
			for _, keyless := range point.Resource.Keyless {
				attestorSet.Entries = append(attestorSet.Entries, kyvernov1.Attestor{
					Keyless: &kyvernov1.KeylessAttestor{
						Issuer:  keyless.Issuer,
						Subject: keyless.Subject,
						Rekor:   &kyvernov1.Rekor{URL: keyless.Url},
					},
				})
			}
		}

		imageVerification.Attestors = append(imageVerification.Attestors, attestorSet)
	}

	return imageVerification
}
