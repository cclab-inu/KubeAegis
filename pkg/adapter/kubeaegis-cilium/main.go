package main

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	pb "github.com/cclab-inu/KubeAegis/api/grpc"
	"github.com/cclab-inu/KubeAegis/pkg/adapter/k8s"
	"github.com/cclab-inu/KubeAegis/pkg/adapter/kubeaegis-cilium/manager"
)

const adapterName = "kubeaegis-cilium"

type server struct {
	pb.UnimplementedPolicyServiceServer
}

func (s *server) DispatchPolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.PolicyResponse, error) {
	logger := ctrl.Log.WithName("main")
	logger.Info("KubeAegis arrived", "KubeAegis.Name", in.GetPolicyName(), "KubeAegis.Namespace", in.GetPolicyNamespace())

	kspname, _ := manager.Run(ctx, logger, in.GetPolicyName(), in.GetPolicyNamespace())

	return &pb.PolicyResponse{
		Success:           true,
		Message:           in.GetPolicyName(),
		AdapterPolicyName: kspname,
	}, nil
}

func (s *server) NotifyPolicyDeletion(ctx context.Context, in *pb.PolicyDeletionRequest) (*pb.PolicyDeletionResponse, error) {
	logger := ctrl.Log.WithName("main")
	logger.Info("CiliumPolicy deleted", "cilium.Name", in.GetPolicyName(), "cilium.Namespace", in.GetPolicyNamespace())

	return &pb.PolicyDeletionResponse{
		Success: true,
		Message: "CiliumPolicy deletion processed successfully",
	}, nil
}

func main() {
	ctrl.SetLogger(zap.New())
	logger := ctrl.Log.WithName("main")

	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	k8sClient := k8s.NewOrDie(scheme)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	ctrl.LoggerInto(ctx, logger)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		logger.Info("Shutdown signal received, exiting...")
		updateAdapterStatus(ctx, k8sClient, logger, "offline")
		cancelFunc()
		os.Exit(1)
	}()
	logger.Info("Cilium adapter started")
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		logger.Error(err, "failed to listen on port 50052")
		os.Exit(1)
	}
	updateAdapterStatus(ctx, k8sClient, logger, "online")
	//go kspwatcher.WatchKsps(ctx, logger)

	s := grpc.NewServer()
	pb.RegisterPolicyServiceServer(s, &server{})

	logger.Info("gRPC server listening on port 50052")
	if err := s.Serve(lis); err != nil {
		logger.Error(err, "failed to serve gRPC server")
		os.Exit(1)
	}
}

func updateAdapterStatus(ctx context.Context, k8sClient client.Client, logger logr.Logger, status string) {
	var configMap corev1.ConfigMap
	err := k8sClient.Get(ctx, types.NamespacedName{Name: "adapter-config", Namespace: "default"}, &configMap)
	if err != nil {
		logger.Error(err, "failed to get ConfigMap")
		return
	}

	var configData map[string]interface{}
	err = json.Unmarshal([]byte(configMap.Data["config"]), &configData)
	if err != nil {
		logger.Error(err, "failed to unmarshal ConfigMap data")
		return
	}

	if adapterConfig, ok := configData[adapterName].(map[string]interface{}); ok {
		adapterConfig["status"] = status
	} else {
		logger.Error(nil, "adapter not found in ConfigMap", "adapterName", adapterName)
		return
	}

	updatedConfigData, err := json.Marshal(configData)
	if err != nil {
		logger.Error(err, "failed to marshal updated ConfigMap data")
		return
	}

	configMap.Data["config"] = string(updatedConfigData)
	err = k8sClient.Update(ctx, &configMap)
	if err != nil {
		logger.Error(err, "failed to update ConfigMap")
	} else {
		logger.Info("Adapter status updated", "adapterName", adapterName, "status", status)
	}
}
