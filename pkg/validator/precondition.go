// Precondition Validator

package validator

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	processor "github.com/cclab-inu/KubeAegis/pkg/adapter/processor"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ValidatePrecondition(ctx context.Context, k8sClient client.Client, kap *v1.KubeAegisPolicy) []error {
	var errs []error

	for _, intentRequest := range kap.Spec.IntentRequest {
		switch intentRequest.Type {
		case "network":
			if err := validateNetworkIntentRequest(ctx, k8sClient, intentRequest); err != nil {
				errs = append(errs, err)
			}
		case "system":
			if err := validateSystemIntentRequest(ctx, k8sClient, kap); err != nil {
				errs = append(errs, err)
			}
		case "cluster":
			if err := validateClusterIntentRequest(ctx, k8sClient, intentRequest); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func validateNetworkIntentRequest(ctx context.Context, k8sClient client.Client, intentRequest v1.IntentRequest) error {
	var matchLabels map[string]string
	var err error
	if len(intentRequest.Selector.CEL) > 0 {
		namespace := ""
		matchLabels, err = processor.ProcessCEL(ctx, k8sClient, namespace, intentRequest.Selector.CEL)
		if err != nil {
			return err
		}
	} else if len(intentRequest.Selector.Match) > 0 {
		matchLabels, err = processor.ProcessMatchLabels(intentRequest.Selector.Match)
		if err != nil {
			return err
		}
	}

	if err := validatePortListening(ctx, k8sClient, intentRequest, matchLabels); err != nil {
		return err
	}

	if err := validateCIDR(intentRequest); err != nil {
		return err
	}

	return nil
}

func validateSystemIntentRequest(ctx context.Context, k8sClient client.Client, kap *v1.KubeAegisPolicy) error {
	for _, intent := range kap.Spec.IntentRequest {
		if err := validateExecutableScriptsAndCommands(intent); err != nil {
			return err
		}

		if err := validateExecutableFilesAndDirectories(intent); err != nil {
			return err
		}

		if err := validateSystemCalls(intent); err != nil {
			return err
		}

	}

	return nil

}

func validateClusterIntentRequest(ctx context.Context, k8sClient client.Client, intentRequest v1.IntentRequest) error {
	var errs []error

	for _, point := range intentRequest.Rule.ActionPoint {
		matchLabels := intentRequest.Selector.Match[0].MatchLabels
		namespace := intentRequest.Selector.Match[0].Namespace
		for _, detailMap := range point.Resource.Details {
			switch point.Resource.Kind {
			case "annotations":
				errs = append(errs, validatePodAnnotations(ctx, k8sClient, matchLabels, namespace, detailMap)...)
			case "label":
				errs = append(errs, validatePodLabels(ctx, k8sClient, matchLabels, namespace, detailMap)...)
			}
		}
	}
	if err := validateImages(intentRequest); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("found errors: %v", errs)
	}

	return nil
}

func validatePortListening(ctx context.Context, k8sClient client.Client, intentRequest v1.IntentRequest, matchLabels map[string]string) error {
	var namespace string
	if len(intentRequest.Selector.Match) > 0 {
		namespace = intentRequest.Selector.Match[0].Namespace
	} else {
		namespace = "default"
	}

	if len(intentRequest.Selector.Match) == 0 && len(intentRequest.Selector.CEL) == 0 {
		return errors.New("no matches found in the selector")
	}
	rules := append(intentRequest.Rule.To, intentRequest.Rule.From...)

	// Check if netPolDetails slice is empty
	if len(rules) == 0 {
		return errors.New("rules slice is empty")
	}

	for _, rule := range rules {
		if rule.Port == "" {
			return errors.New("port is empty")
		}

		expectedPort, err := strconv.ParseInt(rule.Port, 10, 32)
		if err != nil {
			return errors.Errorf("error parsing expected port: %v", err)
		}

		var pods corev1.PodList
		listOpts := []client.ListOption{
			client.InNamespace(namespace),
			client.MatchingLabels(matchLabels),
		}
		if err := k8sClient.List(ctx, &pods, listOpts...); err != nil {
			return errors.Errorf("error fetching pods: %v", err)
		}
		if len(pods.Items) == 0 {
			return errors.Errorf("no pods found matching the selector in namespace %s", namespace)
		}

		portFound := false
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				for _, port := range container.Ports {
					if port.ContainerPort == int32(expectedPort) && strings.EqualFold(string(port.Protocol), string(rule.Protocol)) {
						portFound = true
						break
					}
				}
				if portFound {
					break
				}
			}
			if !portFound {
				return errors.Errorf("no containers found in namespace %s with labels %v listening on the expected port %d with protocol %s", namespace, matchLabels, expectedPort, string(rule.Protocol))
			}
		}
		if err := validateProtocol(string(rule.Protocol)); err != nil {
			return err
		}
	}

	return nil
}

func validateProtocol(protocol string) error {
	validProtocols := []string{"TCP", "UDP", "ICMP"}
	for _, validProtocol := range validProtocols {
		if protocol == validProtocol {
			return nil
		}
	}
	return errors.Errorf("invalid protocol: %s. Must be one of %v", protocol, validProtocols)
}

func validateCIDR(intentRequest v1.IntentRequest) error {
	for _, rule := range append(intentRequest.Rule.To, intentRequest.Rule.From...) {
		for _, cidr := range rule.Labels {
			if _, _, err := net.ParseCIDR(cidr); err != nil {
				return errors.New("invalid CIDR: " + cidr)
			}
		}
	}

	return nil
}

func validateSystemCalls(intentRequest v1.IntentRequest) error {
	knownSystemCalls := []string{"open", "read", "write", "close", "unlink", "rmdir"}
	for _, point := range intentRequest.Rule.ActionPoint {
		if point.SubType == "kprobe" || point.SubType == "tracepoint" {
			syscalls := strings.Split(point.Resource.Syscall, ",")
			for _, call := range syscalls {
				if !contains(knownSystemCalls, strings.TrimSpace(call)) {
					return errors.Errorf("invalid system call: %s", call)
				}
			}
		}
	}

	return nil
}

func validateExecutableFilesAndDirectories(intentRequest v1.IntentRequest) error {
	for _, point := range intentRequest.Rule.ActionPoint {
		for _, path := range point.Resource.Path {
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				return errors.Errorf("file or directory does not exist: %s", path)
			}
		}
	}
	return nil
}

func validateExecutableScriptsAndCommands(intentRequest v1.IntentRequest) error {
	for _, point := range intentRequest.Rule.ActionPoint {
		if point.SubType == "uprobes" {
			if point.Resource.Symbol == "" {
				return errors.Errorf("missing required symbol for uprobes point")
			}
		}
	}
	return nil
}

// validateFileExistenceAndReadOnly checks if the specified files exist and are set to read-only.
func validateFileExistenceAndReadOnly(intentRequest v1.IntentRequest) error {
	for _, point := range intentRequest.Rule.ActionPoint {
		for _, filePath := range point.Resource.Path {
			// 파일 존재 여부 확인
			fileInfo, err := os.Stat(filePath)
			if errors.Is(err, os.ErrNotExist) {
				return errors.Errorf("file does not exist: %s", filePath)
			}

			if !fileInfo.Mode().IsRegular() || (fileInfo.Mode().Perm()&0222) != 0 {
				return errors.Errorf("file is not read-only: %s", filePath)
			}
		}
	}
	return nil
}

// RunCMD executes a given command and returns its output.
func RunCMD(cmd *exec.Cmd) error {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "command failed, check your command")
	}

	fmt.Printf("Running command: %s\nOutput: %s\n", strings.Join(cmd.Args, " "), string(out))
	return nil
}

func validateImages(intentRequest v1.IntentRequest) error {
	for _, point := range intentRequest.Rule.ActionPoint {
		if point.SubType == "verifyImage" {
			if err := verifyImage(point.Resource); err != nil {
				return err
			}
		}
	}
	return nil
}

func verifyImage(resource v1.EventMatchResource) error {
	for _, detailMap := range resource.Details {
		for image, _ := range detailMap {
			resp, err := http.Get("https://index.docker.io/v2/" + image)
			if err != nil {
				return errors.Errorf("error checking image existence: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return errors.Errorf("image %s does not exist", image)
			}
		}
	}
	return nil
}

func validatePodAnnotations(ctx context.Context, k8sClient client.Client, matchLabels map[string]string, namespace string, details map[string]string) []error {
	var errs []error
	var pods corev1.PodList
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(matchLabels),
	}
	if err := k8sClient.List(ctx, &pods, listOpts...); err != nil {
		errs = append(errs, fmt.Errorf("error fetching pods: %v", err))
		return errs
	}
	for _, pod := range pods.Items {
		for key, expectedValue := range details {
			if value, ok := pod.Annotations[key]; ok {
				if value != expectedValue {
					errs = append(errs, fmt.Errorf("pod %s in namespace %s does not comply with the annotation %s: expected %s, got %s", pod.Name, pod.Namespace, key, expectedValue, value))
				}
			}
		}
	}
	return errs
}

func validatePodLabels(ctx context.Context, k8sClient client.Client, matchLabels map[string]string, namespace string, labels map[string]string) []error {
	var errs []error
	var pods corev1.PodList
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(matchLabels),
	}
	if err := k8sClient.List(ctx, &pods, listOpts...); err != nil {
		errs = append(errs, fmt.Errorf("error fetching pods: %v", err))
		return errs
	}

	for _, pod := range pods.Items {
		for key, expectedValue := range labels {
			if value, ok := pod.Labels[key]; ok {
				if value != expectedValue {
					errs = append(errs, fmt.Errorf("pod %s in namespace %s does not comply with the label %s: expected %s, got %s", pod.Name, pod.Namespace, key, expectedValue, value))
				}
			} else {
				errs = append(errs, fmt.Errorf("pod %s in namespace %s is missing required label %s", pod.Name, pod.Namespace, key))
			}
		}
	}
	return errs
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
