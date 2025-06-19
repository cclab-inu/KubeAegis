package reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
)

type AdapterConfig struct {
	Address        string              `json:"address"`
	SupportedTypes map[string][]string `json:"supportedTypes"`
	Status         string              `json:"status"`
}

// Report structure to hold the reporting data
type Report struct {
	PolicyName        string    `json:"kapName"`
	Namespace         string    `json:"kapNamespace"`
	PolicyType        string    `json:"policyType"`
	PolicyActions     string    `json:"policyActions"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
	ValidationErrors  int32     `json:"validationErrors"`
	AdapterName       string    `json:"adapterName"`
	AdapterStatus     string    `json:"adapterStatus"`
	Policies          string    `json:"adapterPolicies"`
}

// GenerateReport collects data and creates a report in JSON format
func GenerateReport(ctx context.Context, logger logr.Logger, kap *v1.KubeAegisPolicy, configMap corev1.ConfigMap, adapterPolicy string) error {
	if !kap.Spec.EnableReporting {
		return nil // Reporting is not enabled, skip report generation
	}

	// Initialize reportData with known values from kap
	reportData := &Report{
		PolicyName:        kap.Name,
		Namespace:         kap.Namespace,
		PolicyType:        "",
		PolicyActions:     "",
		CreationTimestamp: kap.CreationTimestamp.Time,
		ValidationErrors:  kap.Status.NumberOfAPs,
		AdapterName:       "",
		AdapterStatus:     "",
		Policies:          adapterPolicy,
	}
	for _, intentRequest := range kap.Spec.IntentRequest {
		reportData.PolicyType = intentRequest.Type
		reportData.PolicyActions = intentRequest.Rule.Action
	}

	adapterConfigs := map[string]AdapterConfig{}
	if err := json.Unmarshal([]byte(configMap.Data["config"]), &adapterConfigs); err != nil {
		return err
	}

	// Iterate over the intentRequests and dispatch them to the supported adapters.
	for _, intentRequest := range kap.Spec.IntentRequest {
		subType := intentRequest.Rule.ActionPoint[0].SubType
		adaptersToNotify := getSupportedAdapters(adapterConfigs, intentRequest.Type, subType)
		for _, adapterName := range adaptersToNotify {
			adapterConfig, ok := adapterConfigs[adapterName]
			if !ok {
				continue
			}
			reportData.AdapterName = adapterName
			reportData.AdapterStatus = adapterConfig.Status
		}
	}

	// Serialize the report data to JSON
	reportBytes, err := json.MarshalIndent(reportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report data: %v", err)
	}

	// Write the report to a file
	reportDir := "./report"
	if err := os.MkdirAll(reportDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create report directory: %v", err)
	}

	reportPath := filepath.Join(reportDir, fmt.Sprintf("%s-%s.json", kap.Name, kap.Namespace))
	if err := os.WriteFile(reportPath, reportBytes, 0644); err != nil {
		return fmt.Errorf("failed to write report file: %v", err)
	}

	logger.Info("Report created", "KubeAegis.Name", kap.Name, "KubeAegis.Namespace", kap.Namespace)
	return nil
}

func getSupportedAdapters(adapterConfigs map[string]AdapterConfig, intentType, subType string) []string {
	var supportedAdapters []string
	for adapterName, config := range adapterConfigs {
		supportedTypes, ok := config.SupportedTypes[intentType]
		if !ok {
			continue
		}

		for _, sType := range supportedTypes {
			if sType == subType {
				supportedAdapters = append(supportedAdapters, adapterName)
				break
			}
		}
	}
	if len(supportedAdapters) == 0 {
		return nil
	}
	return supportedAdapters
}
