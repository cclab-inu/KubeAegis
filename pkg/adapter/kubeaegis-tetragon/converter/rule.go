package converter

import v1 "github.com/cclab-inu/KubeAegis/api/v1"

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
