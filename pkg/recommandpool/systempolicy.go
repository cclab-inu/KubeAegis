package recommandpool

import (
	karmorv1 "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/api/security.kubearmor.com/v1"
)

// ----------------------------
// KubeArmor (11)
// ----------------------------

// ----------------------------
// KubeArmor Policy: Default Values
// ----------------------------

// SetDefaultValues sets default values for a KubeArmor policy
func SetDefaultValues(ksp *karmorv1.KubeArmorPolicy) {
	ruleDescription = "This function sets default values for a KubeArmor policy. It ensures that if no specific network protocols or capabilities are matched, the policy defaults to allowing 'raw' network protocol and 'lease' capability."
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

// ----------------------------
// KubeArmor Policy: Process
// ----------------------------

// HandleProcessPath adds process path matches to a KubeArmor policy
func HandleProcessPath(ksp *karmorv1.KubeArmorPolicy, path string) {
	ruleDescription = "This function adds process path matches to a KubeArmor policy. It adds specified file paths that processes can access to the policy, helping control which paths can be used by processes."
	ksp.Spec.Process.MatchPaths = append(ksp.Spec.Process.MatchPaths, karmorv1.ProcessPathType{
		Path: karmorv1.MatchPathType(path),
	})
}

// HandleProcessPattern adds process pattern matches to a KubeArmor policy
func HandleProcessPattern(ksp *karmorv1.KubeArmorPolicy, pattern string) {
	ruleDescription = "This function adds process pattern matches to a KubeArmor policy. It includes specified patterns for process names, allowing the policy to control which process names are permissible."
	ksp.Spec.Process.MatchPatterns = append(ksp.Spec.Process.MatchPatterns, karmorv1.ProcessPatternType{
		Pattern: pattern,
	})
}

// HandleProcessDirectory adds process directory matches to a KubeArmor policy
func HandleProcessDirectory(ksp *karmorv1.KubeArmorPolicy, dir string) {
	ruleDescription = "This function adds process directory matches to a KubeArmor policy. It specifies directories that processes are allowed to access, ensuring directory-level control over process execution."
	ksp.Spec.Process.MatchDirectories = append(ksp.Spec.Process.MatchDirectories, karmorv1.ProcessDirectoryType{
		Directory: karmorv1.MatchDirectoryType(dir),
	})
}

// ----------------------------
// KubeArmor Policy: File
// ----------------------------

// HandleFilePath adds file path matches to a KubeArmor policy
func HandleFilePath(ksp *karmorv1.KubeArmorPolicy, path string) {
	ruleDescription = "This function adds file path matches to a KubeArmor policy. It allows the policy to specify which file paths are accessible, providing fine-grained control over file access."
	ksp.Spec.File.MatchPaths = append(ksp.Spec.File.MatchPaths, karmorv1.FilePathType{
		Path: karmorv1.MatchPathType(path),
	})
}

// HandleFilePattern adds file pattern matches to a KubeArmor policy
func HandleFilePattern(ksp *karmorv1.KubeArmorPolicy, pattern string) {
	ruleDescription = "This function adds file pattern matches to a KubeArmor policy. It specifies patterns for file names, ensuring that only files matching the patterns can be accessed."
	ksp.Spec.File.MatchPatterns = append(ksp.Spec.File.MatchPatterns, karmorv1.FilePatternType{
		Pattern: pattern,
	})
}

// HandleFileDirectory adds file directory matches to a KubeArmor policy
func HandleFileDirectory(ksp *karmorv1.KubeArmorPolicy, dir string) {
	ruleDescription = "This function adds file directory matches to a KubeArmor policy. It specifies directories that files can be accessed from, helping to control file access at the directory level."
	ksp.Spec.File.MatchDirectories = append(ksp.Spec.File.MatchDirectories, karmorv1.FileDirectoryType{
		Directory: karmorv1.MatchDirectoryType(dir),
	})
}

// ----------------------------
// KubeArmor Policy: Network
// ----------------------------

// HandleNetworkProtocol adds network protocol matches to a KubeArmor policy
func HandleNetworkProtocol(ksp *karmorv1.KubeArmorPolicy, protocol string) {
	ruleDescription = "This function adds network protocol matches to a KubeArmor policy. It allows the policy to specify which network protocols are allowed, enhancing control over network traffic."
	ksp.Spec.Network.MatchProtocols = append(ksp.Spec.Network.MatchProtocols, karmorv1.MatchNetworkProtocolType{
		Protocol: karmorv1.MatchNetworkProtocolStringType(protocol),
	})
}

// ----------------------------
// KubeArmor Policy: Capabilities
// ----------------------------

// HandleCapabilities adds capability matches to a KubeArmor policy
func HandleCapabilities(ksp *karmorv1.KubeArmorPolicy, capability string) {
	ruleDescription = "This function adds capability matches to a KubeArmor policy. It specifies Linux capabilities that processes can have, ensuring that only allowed capabilities are granted."
	ksp.Spec.Capabilities.MatchCapabilities = append(ksp.Spec.Capabilities.MatchCapabilities, karmorv1.MatchCapabilitiesType{
		Capability: karmorv1.MatchCapabilitiesStringType(capability),
	})
}

// ----------------------------
// KubeArmor Policy: Syscalls
// ----------------------------

// HandleSyscallPath adds syscall path matches to a KubeArmor policy
func HandleSyscallPath(ksp *karmorv1.KubeArmorPolicy, path string) {
	ruleDescription = "This function adds syscall path matches to a KubeArmor policy. It specifies paths for system calls, helping to control which system calls can be made from which paths."
	ksp.Spec.Syscalls.MatchPaths = append(ksp.Spec.Syscalls.MatchPaths, karmorv1.SyscallMatchPathType{
		Path: karmorv1.MatchSyscallPathType(path),
	})
}

// HandleSyscall adds syscall matches to a KubeArmor policy
func HandleSyscall(ksp *karmorv1.KubeArmorPolicy, syscall string) {
	ruleDescription = "This function adds syscall matches to a KubeArmor policy. It allows the policy to specify which system calls are permissible, controlling system-level interactions."
	ksp.Spec.Syscalls.MatchSyscalls = append(ksp.Spec.Syscalls.MatchSyscalls, karmorv1.SyscallMatchType{
		Syscalls: []karmorv1.Syscall{karmorv1.Syscall(syscall)},
	})
}

// ================================================================================================================
// ----------------------------
// Tetragon
// ----------------------------
// ================================================================================================================

/*

// ConvertContainerSelector converts KubeAegisPolicy containerSelector to Tetragon containerSelector
func ConvertContainerSelector(selector v1.ActionPoint) api.ContainerSelector {
	return api.ContainerSelector{
		MatchExpressions: convertMatchExpressions(selector),
		MatchLabels:      selector.Resource.Path,
	}
}

func convertMatchExpressions(match []v1.EventMatchResource) []api.MatchExpression {
	var matchExpressions []api.MatchExpression
	for _, m := range match {
		expr := api.MatchExpression{
			Key:      m.Keys,
			Operator: m.Details,
			Values:   m.Methods,
		}
		matchExpressions = append(matchExpressions, expr)
	}
	return matchExpressions
}

// ConvertEnforcers converts KubeAegisPolicy enforcers to Tetragon enforcers
func ConvertEnforcers(rules []v1.Rule) []api.Enforcer {
	var enforcers []api.Enforcer
	for _, r := range rules {
		enforcer := api.Enforcer{
			Calls: r.ActionPoint.Methods, // Adjust as necessary
		}
		enforcers = append(enforcers, enforcer)
	}
	return enforcers
}

// ConvertKprobes converts KubeAegisPolicy kprobes to Tetragon kprobes
func ConvertKprobes(rules []v1.Rule) []api.Kprobe {
	var kprobes []api.Kprobe
	for _, r := range rules {
		kprobe := api.Kprobe{
			Call:      r.ActionPoint[0].Resource.Path, // Adjust as necessary
			Args:      convertKprobeArgs(r.ActionPoint[0].Resource.Args),
			Selectors: convertKprobeSelectors(r.ActionPoint),
			Tags:      r.ActionPoint[0].Resource.Keys,
		}
		kprobes = append(kprobes, kprobe)
	}
	return kprobes
}

func convertKprobeArgs(args []string) []api.KprobeArg {
	var kprobeArgs []api.KprobeArg
	for _, arg := range args {
		kprobeArg := api.KprobeArg{
			Index: 0, // Set appropriate index
			Type:  arg,
		}
		kprobeArgs = append(kprobeArgs, kprobeArg)
	}
	return kprobeArgs
}

func convertKprobeSelectors(filters []v1.ActionPoint) []api.KprobeSelector {
	var selectors []api.KprobeSelector
	for _, f := range filters {
		selector := api.KprobeSelector{
			MatchActions: convertMatchActions(f.Resource),
			MatchArgs:    convertMatchArgs(f.Resource),
		}
		selectors = append(selectors, selector)
	}
	return selectors
}

func convertMatchActions(actions []v1.EventMatchResource) []api.MatchAction {
	var matchActions []api.MatchAction
	for _, a := range actions {
		action := api.MatchAction{
			Action: a.Args,
		}
		matchActions = append(matchActions, action)
	}
	return matchActions
}

func convertMatchArgs(args []v1.EventMatchResource) []api.MatchArg {
	var matchArgs []api.MatchArg
	for _, arg := range args {
		matchArg := api.MatchArg{
			Index:    arg.Args,
			Operator: arg.Kind,
			Values:   arg.Details,
		}
		matchArgs = append(matchArgs, matchArg)
	}
	return matchArgs
}

// ConvertPodSelector converts KubeAegisPolicy podSelector to Tetragon podSelector
func ConvertPodSelector(selector v1.Selector) api.PodSelector {
	return api.PodSelector{
		MatchExpressions: convertMatchExpressions(selector.Match),
		MatchLabels:      selector.Match,
	}
}

// ConvertLists converts KubeAegisPolicy lists to Tetragon lists
func ConvertLists(lists []v1.ActionPoint) []api.ListSpec {
	var listSpecs []api.ListSpec
	for _, l := range lists {
		listSpec := api.ListSpec{
			Name:      l.Resource.Name,
			Pattern:   l.Resource.Pattern,
			Type:      l.Resource.Args,
			Validated: l.Resource.List,
			Values:    l.Resource.List,
		}
		listSpecs = append(listSpecs, listSpec)
	}
	return listSpecs
}

*/
