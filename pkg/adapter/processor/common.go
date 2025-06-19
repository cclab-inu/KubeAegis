package processor

import (
	"regexp"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func checkNegation(expr string) bool {
	isNegated := strings.HasPrefix(expr, "'!") || strings.HasPrefix(expr, "!") || strings.HasPrefix(expr, "!='") || strings.Contains(expr, " != ")
	return isNegated
}

func preprocessExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	expr = regexp.MustCompile(`^['"]|['"]$`).ReplaceAllString(expr, "")
	expr = strings.ReplaceAll(expr, `\"`, `"`)
	expr = strings.ReplaceAll(expr, `\'`, `'`)
	expr = strings.Replace(expr, `\"`, `"`, -1)
	expr = regexp.MustCompile(`^['"]|['"]$`).ReplaceAllString(expr, "")
	if strings.Count(expr, "\"")%2 != 0 {
		expr += "\""
	} else if strings.Count(expr, "'")%2 != 0 {
		expr += "'"
	}

	return expr
}

func extractLabelsFromExpression(expr string, podList corev1.PodList, isNegated bool) map[string]string {
	labels := make(map[string]string)

	if strings.Contains(expr, "==") || strings.Contains(expr, "!=") {
		key, value := parseKeyValueExpression(expr)
		labels[key] = value
	} else if strings.Contains(expr, ".contains(") {
		key, value := parseFunctionExpression(expr, "contains")
		if key != "" && value != "" {
			labels[key] = value
		}
	} else if strings.Contains(expr, " in ") {
		key, values := parseInExpression(expr)
		for _, pod := range podList.Items {
			labelValue, exists := pod.Labels[key]
			if !exists {
				continue
			}
			if contains(values, labelValue) {
				labels[key] = labelValue
			}
		}
	} else if strings.Contains(expr, ".startsWith(") {
		labels = parseStartsWithEndsWithExpression(expr, podList, "startsWith")
	} else if strings.Contains(expr, ".endsWith(") {
		labels = parseStartsWithEndsWithExpression(expr, podList, "endsWith")
	} else if strings.Contains(expr, ".matches(") {
		labels = parseMatchesExpression(expr, podList)
	}
	if isNegated {
		labels = excludeLabels(podList, labels)
	}
	return labels
}

func parseKeyValueExpression(expr string) (string, string) {
	expr = preprocessExpression(expr)
	var operator string
	if strings.Contains(expr, "==") {
		operator = "=="
	} else if strings.Contains(expr, "!=") {
		operator = "!="
	} else {
		return "", ""
	}

	parts := strings.SplitN(expr, operator, 2)
	if len(parts) != 2 {
		return "", ""
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	key = strings.TrimPrefix(key, "labels[")
	key = strings.TrimSuffix(key, "]")
	key = strings.Trim(key, `"'`)
	value = strings.Trim(value, `"'`)
	return key, value
}

func parseFunctionExpression(expr string, functionName string) (string, string) {
	start := strings.Index(expr, `labels["`) + len(`labels["`)
	if start == -1 {
		return "", ""
	}
	end := strings.Index(expr[start:], `"]`)
	if end == -1 {
		return "", ""
	}
	key := expr[start : start+end]

	functionStart := strings.Index(expr, functionName+"(\"") + len(functionName+"(\"")
	functionEnd := strings.LastIndex(expr, "\")")
	if functionStart == -1 || functionEnd == -1 || functionStart >= functionEnd {
		return "", ""
	}
	value := expr[functionStart:functionEnd]

	return key, value
}

func parseInExpression(expr string) (string, []string) {
	start := strings.Index(expr, `labels["`) + len(`labels["`)
	if start == -1 {
		return "", nil
	}
	end := strings.Index(expr[start:], `"]`)
	if end == -1 {
		return "", nil
	}
	key := expr[start : start+end]

	valuesStart := strings.Index(expr, " in [") + len(" in [")
	valuesEnd := strings.LastIndex(expr, "]")
	if valuesStart == -1 || valuesEnd == -1 || valuesStart >= valuesEnd {
		return "", nil
	}
	valuesString := expr[valuesStart:valuesEnd]
	valuesParts := strings.Split(valuesString, ",")

	var values []string
	for _, part := range valuesParts {
		value := strings.TrimSpace(part)
		value = strings.Trim(value, "\"'")
		values = append(values, value)
	}

	return key, values
}

func parseStartsWithEndsWithExpression(expr string, podList corev1.PodList, functionName string) map[string]string {
	labels := make(map[string]string)
	key, pattern := parseFunctionExpression(expr, functionName)

	for _, pod := range podList.Items {
		labelValue, exists := pod.Labels[key]
		if !exists {
			continue
		}

		var match bool
		if functionName == "startsWith" && strings.HasPrefix(labelValue, pattern) {
			match = true
		} else if functionName == "endsWith" && strings.HasSuffix(labelValue, pattern) {
			match = true
		}

		if match {
			labels[key] = labelValue
		}
	}

	return labels
}

func parseMatchesExpression(expr string, podList corev1.PodList) map[string]string {
	key, pattern := parseFunctionExpression(expr, "matches")
	labels := make(map[string]string)

	regex, _ := regexp.Compile(pattern)

	for _, pod := range podList.Items {
		labelValue, exists := pod.Labels[key]
		if !exists {
			continue
		}

		if regex.MatchString(labelValue) {
			labels[key] = labelValue
		}
	}

	return labels
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func excludeLabels(podList corev1.PodList, excludeMap map[string]string) map[string]string {
	remainingLabels := make(map[string]string)

	for _, pod := range podList.Items {
		exclude := false
		for excludeKey, excludeValue := range excludeMap {
			podLabelValue, exists := pod.Labels[excludeKey]
			if exists && podLabelValue == excludeValue {
				exclude = true
				break
			}
		}

		if !exclude {
			for labelKey, labelValue := range pod.Labels {
				if labelKey != "pod-template-hash" {
					remainingLabels[labelKey] = labelValue
				}
			}
		}
	}

	return remainingLabels
}
