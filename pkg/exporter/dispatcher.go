package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pb "github.com/cclab-inu/KubeAegis/api/grpc"
	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	"github.com/cclab-inu/KubeAegis/pkg/statusmanager"
)

// AdapterConfig represents the configuration for a specific adapter.
type AdapterConfig struct {
	Address        string              `json:"address"`
	SupportedTypes map[string][]string `json:"supportedTypes"`
	Status         string              `json:"status"`
}

// DispatchPolicyToAdapters sends the policy to the appropriate adapters based on the type and subtype.
func DispatchPolicyToAdapters(ctx context.Context, k8sClient client.Client, logger logr.Logger, kap *v1.KubeAegisPolicy, configMap corev1.ConfigMap) error {
	// Parse the adapter configurations.
	adapterConfigs := map[string]AdapterConfig{}
	if err := json.Unmarshal([]byte(configMap.Data["config"]), &adapterConfigs); err != nil {
		return err
	}

	// Iterate over the intentRequests and dispatch them to the supported adapters.
	for _, intentRequest := range kap.Spec.IntentRequest {
		var adaptersToNotify []string
		var subType string

		switch intentRequest.Type {
		case "system", "cluster":
			// System type handling
			subType = intentRequest.Rule.ActionPoint[0].SubType
			adaptersToNotify = GetSupportedAdapters(adapterConfigs, logger, intentRequest.Type, subType)
		case "network":
			// Network type handling. Ensure From and To slices are not empty before accessing.
			if len(intentRequest.Rule.From) > 0 {
				subType = intentRequest.Rule.From[0].Kind
				adaptersToNotify = GetSupportedAdapters(adapterConfigs, logger, intentRequest.Type, subType)
			} else if len(intentRequest.Rule.To) > 0 {
				subType = intentRequest.Rule.To[0].Kind
				adaptersToNotify = GetSupportedAdapters(adapterConfigs, logger, intentRequest.Type, subType)
			}
		}
		for _, adapterName := range adaptersToNotify {
			adapterConfig, exists := adapterConfigs[adapterName]
			if !exists || adapterConfig.Status == "offline" {
				logger.Info("Adapter is offline, will retry", "Adapter.Name", adapterName)
				go retryDispatchPolicy(ctx, k8sClient, logger, kap, adapterName) // Execute retry logic asynchronously
				continue
			}

			response, err := DispatchPolicy(ctx, adapterConfig.Address, kap)
			if err != nil {
				logger.Error(err, "error sending policy to adapter", "Adapter.Name", adapterName)
			} else {
				logger.Info("Policy dispatched to adapter", "Adapter.Name", adapterName)
			}
			if response.Success {
				adapterPolicy := response.AdapterPolicyName
				if err := statusmanager.NotifyReporter(ctx, kap, configMap, adapterPolicy); err != nil {
					logger.Error(err, "failed to notify reporter after policy dispatch")
				}
			}
		}
	}

	return nil
}

// getSupportedAdapters returns a slice of adapter names that support the given type and subtype.
func GetSupportedAdapters(adapterConfigs map[string]AdapterConfig, logger logr.Logger, intentType, subType string) []string {
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

	logger.Info("Adapter found", "Adapter.Name", supportedAdapters)
	return supportedAdapters
}

func DispatchPolicy(ctx context.Context, address string, kap *v1.KubeAegisPolicy) (*pb.PolicyResponse, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect gRPC")
	}
	defer conn.Close()

	client := pb.NewPolicyServiceClient(conn)

	req := &pb.PolicyRequest{
		PolicyName:      kap.Name,
		PolicyNamespace: kap.Namespace,
	}

	response, err := client.DispatchPolicy(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send policy")
	}

	return response, nil
}

func retryDispatchPolicy(ctx context.Context, k8sClient client.Client, logger logr.Logger, kap *v1.KubeAegisPolicy, adapterName string) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var configMap corev1.ConfigMap
			err := k8sClient.Get(ctx, types.NamespacedName{Name: "adapter-config", Namespace: "default"}, &configMap)
			if err != nil {
				logger.Error(err, "failed to get ConfigMap for adapter status")
				continue
			}

			var configData map[string]AdapterConfig
			err = json.Unmarshal([]byte(configMap.Data["config"]), &configData)
			if err != nil {
				logger.Error(err, "failed to unmarshal ConfigMap data")
				continue
			}

			adapterConfig, exists := configData[adapterName]
			if !exists || adapterConfig.Status != "online" {
				continue
			}

			_, err = DispatchPolicy(ctx, adapterConfig.Address, kap)
			if err != nil {
				logger.Error(err, "failed to dispatch policy to adapter on retry", "Adapter.Name", adapterName)
			} else {
				logger.Info("Policy dispatched to adapter on retry", "Adapter.Name", adapterName)
				return
			}
		}
	}
}

// NotifyAdapterOfPolicyDeletion notifies all configured adapters about the deletion of a KubeAegisPolicy.
func NotifyAdapterOfPolicyDeletion(ctx context.Context, namespacedName client.ObjectKey, kap *v1.KubeAegisPolicy, configMap corev1.ConfigMap) error {
	var adapterConfigs map[string]AdapterConfig
	err := json.Unmarshal([]byte(configMap.Data["config"]), &adapterConfigs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal adapter config: %w", err)
	}

	for adapterName, config := range adapterConfigs {
		if config.Status != "online" {
			continue // Skip offline adapters
		}

		conn, err := grpc.Dial(config.Address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return fmt.Errorf("failed to connect to adapter %s at %s: %w", adapterName, config.Address, err)
		}
		defer conn.Close()

		client := pb.NewPolicyServiceClient(conn)
		req := &pb.PolicyDeletionRequest{
			PolicyName:      "ksp-" + namespacedName.Name,
			PolicyNamespace: namespacedName.Namespace,
		}

		_, err = client.NotifyPolicyDeletion(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to notify adapter %s of policy deletion: %w", adapterName, err)
		}
	}

	return nil
}
