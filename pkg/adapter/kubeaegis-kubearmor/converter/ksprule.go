package converter

import (
	v1 "github.com/cclab-inu/KubeAegis/api/v1"
	karmorv1 "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1"
)

func setDefaultValues(ksp *karmorv1.KubeArmorPolicy) {
	if len(ksp.Spec.Network.MatchProtocols) == 0 {
		ksp.Spec.Network.MatchProtocols = append(ksp.Spec.Network.MatchProtocols, karmorv1.MatchNetworkProtocolType{
			Protocol: "raw",
		})
	}
	if len(ksp.Spec.Capabilities.MatchCapabilities) == 0 {
		ksp.Spec.Capabilities.MatchCapabilities = append(ksp.Spec.Capabilities.MatchCapabilities, karmorv1.MatchCapabilitiesType{
			Capability: "lease",
		})
	}
}

func handleProcess(ksp *karmorv1.KubeArmorPolicy, point v1.ActionPoint) {
	if len(point.Resource.Path) > 0 {
		for _, path := range point.Resource.Path {
			ksp.Spec.Process.MatchPaths = append(ksp.Spec.Process.MatchPaths, karmorv1.ProcessPathType{
				Path: karmorv1.MatchPathType(path),
			})
		}
	}

	if len(point.Resource.Pattern) > 0 {
		for _, pattern := range point.Resource.Pattern {
			ksp.Spec.Process.MatchPatterns = append(ksp.Spec.Process.MatchPatterns, karmorv1.ProcessPatternType{
				Pattern: pattern,
			})
		}
	}

	if len(point.Resource.Dir) > 0 {
		for _, dir := range point.Resource.Dir {
			ksp.Spec.Process.MatchDirectories = append(ksp.Spec.Process.MatchDirectories, karmorv1.ProcessDirectoryType{
				Directory: karmorv1.MatchDirectoryType(dir),
			})
		}
	}
}

func handleFile(ksp *karmorv1.KubeArmorPolicy, point v1.ActionPoint) {
	if len(point.Resource.Path) > 0 {
		for _, path := range point.Resource.Path {
			ksp.Spec.File.MatchPaths = append(ksp.Spec.File.MatchPaths, karmorv1.FilePathType{
				Path: karmorv1.MatchPathType(path),
			})
		}
	}

	if len(point.Resource.Pattern) > 0 {
		for _, pattern := range point.Resource.Pattern {
			ksp.Spec.File.MatchPatterns = append(ksp.Spec.File.MatchPatterns, karmorv1.FilePatternType{
				Pattern: pattern,
			})
		}
	}

	if len(point.Resource.Dir) > 0 {
		for _, dir := range point.Resource.Dir {
			ksp.Spec.File.MatchDirectories = append(ksp.Spec.File.MatchDirectories, karmorv1.FileDirectoryType{
				Directory: karmorv1.MatchDirectoryType(dir),
			})
		}
	}
}

func handleNetwork(ksp *karmorv1.KubeArmorPolicy, point v1.ActionPoint) {
	if len(point.Resource.Protocol) > 0 {
		if ksp.Spec.Network.MatchProtocols == nil {
			ksp.Spec.Network.MatchProtocols = []karmorv1.MatchNetworkProtocolType{}
		}
		for _, net := range point.Resource.Protocol {
			ksp.Spec.Network.MatchProtocols = append(ksp.Spec.Network.MatchProtocols, karmorv1.MatchNetworkProtocolType{
				Protocol: karmorv1.MatchNetworkProtocolStringType(net),
			})
		}
	}
}

func handleCapabilities(ksp *karmorv1.KubeArmorPolicy, point v1.ActionPoint) {
	if len(point.Resource.Args) > 0 {
		if ksp.Spec.Network.MatchProtocols == nil {
			ksp.Spec.Network.MatchProtocols = []karmorv1.MatchNetworkProtocolType{}
		}
		for _, cap := range point.Resource.Args {
			ksp.Spec.Capabilities.MatchCapabilities = append(ksp.Spec.Capabilities.MatchCapabilities, karmorv1.MatchCapabilitiesType{
				Capability: karmorv1.MatchCapabilitiesStringType(cap),
			})
		}
	}
}

func handleSyscalls(ksp *karmorv1.KubeArmorPolicy, point v1.ActionPoint) {
	if len(point.Resource.Path) > 0 {
		for _, path := range point.Resource.Path {
			ksp.Spec.Syscalls.MatchPaths = append(ksp.Spec.Syscalls.MatchPaths, karmorv1.SyscallMatchPathType{
				Path: karmorv1.MatchSyscallPathType(path),
			})
		}
	}

	if len(point.Resource.Args) > 0 {
		for _, sys := range point.Resource.Args {
			ksp.Spec.Syscalls.MatchSyscalls = append(ksp.Spec.Syscalls.MatchSyscalls, karmorv1.SyscallMatchType{
				Syscalls: []karmorv1.Syscall{karmorv1.Syscall(sys)},
			})
		}
	}
}
