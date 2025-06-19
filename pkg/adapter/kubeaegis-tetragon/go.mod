module github.com/cclab-inu/KubeAegis/pkg/adapter/kubeaegis-tetragon

go 1.24.2

toolchain go1.24.2

require (
	github.com/go-logr/logr v1.4.2
	k8s.io/apimachinery v0.30.1
	k8s.io/kubernetes v1.29.1
	github.com/google/cel-go v0.17.8
	github.com/onsi/ginkgo/v2 v2.19.0
	github.com/onsi/gomega v1.33.1
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.33.0
	k8s.io/api v0.30.0
	k8s.io/pod-security-admission v0.29.0
	sigs.k8s.io/controller-runtime v0.18.4
)

require github.com/gogo/protobuf v1.3.2 // indirect
replace github.com/cclab-inu/KubeAegis => ../..

