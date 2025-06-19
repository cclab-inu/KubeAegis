package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	recommendpool "github.com/cclab-inu/KubeAegis/pkg/recommandpool"
)

type CRD struct {
	Spec struct {
		Group string `json:"group"`
		Names struct {
			Kind       string   `json:"kind"`
			ShortNames []string `json:"shortNames"`
			Singular   string   `json:"singular"`
		} `json:"names"`
		Versions []struct {
			Name   string `json:"name"`
			Schema struct {
				OpenAPIV3Schema struct {
					Properties map[string]Property `json:"properties"`
				} `json:"openAPIV3Schema"`
			} `json:"schema"`
		} `json:"versions"`
	} `json:"spec"`
}

type Property struct {
	Description          string              `json:"description"`
	Type                 string              `json:"type"`
	Properties           map[string]Property `json:"properties,omitempty"`
	Items                *Property           `json:"items,omitempty"`
	AdditionalProperties *Property           `json:"additionalProperties,omitempty"`
	AnyOf                []Property          `json:"anyOf,omitempty"`
}

var importPackages = map[string]string{
	"SampleSpecNamesKind": "importgopkg", // Default package for unknown types
	// network
	"CiliumNetworkPolicy": "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2",
	"GlobalNetworkPolicy": "github.com/projectcalico/api/pkg/apis/projectcalico/v3",
	"NetworkPolicy":       "github.com/projectcalico/api/pkg/apis/projectcalico/v3", // "github.com/projectcalico/libcalico-go/lib/api"
	// system
	"kubearmorpolicies": "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1",
	"TracingPolicy":     "github.com/cilium/tetragon/api/v1/tetragon",
	// cluster
	"ClusterPolicy": "github.com/kyverno/kyverno/api/kyverno/v1",
	"Policy":        "github.com/kyverno/kyverno/api/kyverno/v1",
}

func main() {
	// var name, crdPath, policyType, subtypes string
	var name, crdPath string

	flag.StringVar(&name, "name", os.Getenv("NAME"), "Name of the adapter")
	flag.StringVar(&crdPath, "crd", os.Getenv("CRD"), "Path to the CRD file")
	// flag.StringVar(&policyType, "type", os.Getenv("TYPE"), "Policy type")
	// flag.StringVar(&subtypes, "subtype", os.Getenv("SUBTYPE"), "Subtypes for the policy")
	flag.Parse()

	// if name == "" || crdPath == "" || policyType == "" || subtypes == "" {
	if name == "" || crdPath == "" {
		// fmt.Println("All parameters --name, --crd, --type, and --subtype are required.")
		// fmt.Println("All parameters --name and --crd are required.")
		return
	}

	crdFile, err := os.ReadFile(crdPath)
	if err != nil {
		fmt.Printf("Error reading CRD file: %v\n", err)
		return
	}

	var crd CRD
	if err := json.Unmarshal(crdFile, &crd); err != nil {
		fmt.Printf("Error parsing CRD file: %v\n", err)
		return
	}

	group := crd.Spec.Group
	kind := crd.Spec.Names.Kind
	shortName := kind
	if len(crd.Spec.Names.ShortNames) > 0 {
		shortName = crd.Spec.Names.ShortNames[0]
	} else if crd.Spec.Names.Singular != "" {
		shortName = crd.Spec.Names.Singular
	}
	version := crd.Spec.Versions[0].Name

	fmt.Println("✓ CRD Information: ")
	fmt.Printf("- Group: %s, Kind: %s, Version: %s, ShortName: %s\n", group, kind, version, shortName)
	// fmt.Println()

	adapterName := fmt.Sprintf("kubeaegis-%s", name)
	adapterDir, err := filepath.Abs(fmt.Sprintf("./pkg/adapter/%s", adapterName))
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		return
	}

	// if err := createAdapter(adapterDir, name, group, kind, version, shortName, policyType, subtypes); err != nil {
	if err := createAdapter(adapterDir, name, group, kind, version, shortName); err != nil {
		fmt.Printf("Error creating adapter: %v\n", err)
		return
	}

	port := getNextAvailablePort()
	// if err := updateConfigMap(adapterName, port, policyType, subtypes); err != nil {
	if err := updateConfigMap(adapterName, port); err != nil {
		fmt.Printf("Error updating adapter config: %v\n", err)
		return
	}

	fmt.Println("✓ Adapter created successfully.")

	fieldDescriptions := extractFieldDescriptions(crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties, "")
	apiMethods := getAPIMethodsByType(kind)

	// fmt.Printf("Field Descriptions: %v\n", fieldDescriptions)
	// fmt.Printf("API Methods: %v\n", apiMethods)

	// apiMethods := recommendpool.GetAllAPIMethods()
	recommendedAPIs, err := recommendAPIs(fieldDescriptions, apiMethods)
	// recommendedAPIs, err = recommendNew(fieldDescriptions, apiMethods)
	if err != nil {
		fmt.Printf("Error recommending APIs: %v\n", err)
		return
	}

	fmt.Println("✓ Recommended APIs for each field:")
	for field, apis := range recommendedAPIs {
		fmt.Printf("- %s:\n", field)
		for _, api := range apis {
			fmt.Printf("  - Field: %s, Score: %f\n", api.Field, api.Score)
		}
	}

	if err := saveResultsToFile(recommendedAPIs, name); err != nil {
		fmt.Printf("Error saving results to file: %v\n", err)
		return
	}
}

func saveResultsToFile(recommendedAPIs map[string][]struct {
	Field string
	Score float64
}, name string) error {
	filePath := filepath.Join("pkg", "adapter-maker", fmt.Sprintf("%s.txt", name))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating results file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for api, _ := range recommendedAPIs {
		_, _ = writer.WriteString(fmt.Sprintf("- %s:\n", api))
	}

	return nil
}

// --- Adapter Generation ---
// func createAdapter(adapterDir, name, group, kind, version, shortName, policyType, subtypes string) error {
func createAdapter(adapterDir, name, group, kind, version, shortName string) error {
	templateDir, err := filepath.Abs("./pkg/adapter/kubeaegisple")
	if err != nil {
		return fmt.Errorf("error getting absolute path: %w", err)
	}

	err = copyDir(templateDir, adapterDir)
	if err != nil {
		return fmt.Errorf("error copying template directory: %w", err)
	}

	replacements := map[string]string{
		"SampleSpecGroup":       group,
		"SampleGroupString":     fmt.Sprintf(`"%s"`, group),
		"SampleSpecVersionName": version,
		"SampleVersionString":   fmt.Sprintf(`"%s"`, version),
		"SampleSpecNamesKind":   kind,
		"SampleKindString":      fmt.Sprintf(`"%s"`, kind),
		"SampleResourcePolicy":  kind,
		"sample":                name,
		"short":                 fmt.Sprintf(`"%s"`, shortName),
		"50000":                 getNextAvailablePort(),
	}

	importPkg := getImportPackage(kind)
	replacements["importgopkg"] = importPkg

	err = replacePlaceholders(adapterDir, replacements)
	if err != nil {
		return fmt.Errorf("error replacing placeholders: %w", err)
	}

	return nil
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, src)
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, info.Mode())
	})
}

func replacePlaceholders(dir string, replacements map[string]string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		input, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(input)
		for old, new := range replacements {
			content = strings.ReplaceAll(content, old, new)
		}

		return os.WriteFile(path, []byte(content), info.Mode())
	})

	return err
}

func getNextAvailablePort() string {
	const basePort = 50051
	const maxPort = 60000

	usedPorts := getUsedPorts()

	for port := basePort; port < maxPort; port++ {
		if !usedPorts[port] {
			ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				ln.Close()
				return strconv.Itoa(port)
			}
		}
	}

	fmt.Println("No available port found.")
	os.Exit(1)
	return ""
}

func getUsedPorts() map[int]bool {
	usedPorts := make(map[int]bool)
	configPath := "./pkg/exporter/adapterconfig.yaml"

	file, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Error reading adapter config file: %v\n", err)
		return usedPorts
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "address") {
			parts := strings.Split(line, ":")
			if len(parts) > 2 {
				portStr := strings.TrimSuffix(parts[2], `",`)
				port, err := strconv.Atoi(portStr)
				if err == nil {
					usedPorts[port] = true
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning adapter config file: %v\n", err)
	}

	return usedPorts
}

// func updateConfigMap(name, port, policyType, subtypes string) error {
func updateConfigMap(name, port string) error {
	configPath := "./pkg/exporter/adapterconfig.yaml"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config map file: %w", err)
	}

	newAdapterConfig := fmt.Sprintf(`,
		  "%s": {
			"supportedTypes": {
			  "": [""]
			},
			"address": "localhost:%s",
			"status": "offline"
		  }`, name, port)

	updatedConfig := strings.TrimSuffix(string(data), "\n    }") + newAdapterConfig + "\n    }"
	return os.WriteFile(configPath, []byte(updatedConfig), 0644)
}

func getImportPackage(kind string) string {
	if pkg, exists := importPackages[kind]; exists {
		return pkg
	}
	return importPackages["SampleSpecNamesKind"]
}

func extractFieldDescriptions(properties map[string]Property, path string) map[string]string {
	fieldDescriptions := make(map[string]string)
	for field, prop := range properties {
		fullPath := fmt.Sprintf("%s.%s", path, field)
		fieldDescriptions[fullPath] = prop.Description
		if len(prop.Properties) > 0 {
			nestedDescriptions := extractFieldDescriptions(prop.Properties, fullPath)
			for k, v := range nestedDescriptions {
				fieldDescriptions[k] = v
			}
		}
		if prop.Items != nil && len(prop.Items.Properties) > 0 {
			itemDescriptions := extractFieldDescriptions(prop.Items.Properties, fullPath)
			for k, v := range itemDescriptions {
				fieldDescriptions[k] = v
			}
		}
	}
	return fieldDescriptions
}

func getAPIMethodsByType(kind string) []recommendpool.APIMethod {
	switch kind {
	case "CiliumNetworkPolicy", "GlobalNetworkPolicy", "NetworkPolicy":
		return recommendpool.GetNetworkAPIMethods()
	case "kubearmorpolicies", "TracingPolicy", "KubeArmorPolicy", "TracingPolicyNamespaced", "tracingpolicynamespaced":
		return recommendpool.GetSystemAPIMethods()
	case "ClusterPolicy", "Policy":
		return recommendpool.GetClusterAPIMethods()
	default:
		fmt.Printf("Unknown kind: %s\n", kind)
		return []recommendpool.APIMethod{}
	}
}

func recommendAPIs(fieldDescriptions map[string]string, apiMethods []recommendpool.APIMethod) (map[string][]struct {
	Field string
	Score float64
}, error) {
	inputData := map[string]interface{}{
		"fieldDescriptions": fieldDescriptions,
		"apiMethods":        apiMethods,
	}

	inputJSON, err := json.Marshal(inputData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input data: %w", err)
	}

	scriptPath := filepath.Join("pkg", "adapter-maker", "sRobert_cos_v1.py")

	cmd := exec.Command("python3", scriptPath)
	cmd.Stdin = strings.NewReader(string(inputJSON))
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running Python script: %v\n", err)
		fmt.Printf("Python script stderr: %s\n", errBuf.String())
		return nil, fmt.Errorf("error running Python script: %w", err)
	}

	output := outBuf.String()
	var recommendedAPIs map[string][]struct {
		Field string
		Score float64
	}
	if err := json.Unmarshal([]byte(output), &recommendedAPIs); err != nil {
		fmt.Printf("Error unmarshaling output: %v\n", err)
		fmt.Printf("Python script raw output: %s\n", output)
		return nil, fmt.Errorf("error unmarshaling output: %w", err)
	}

	return recommendedAPIs, nil
}
