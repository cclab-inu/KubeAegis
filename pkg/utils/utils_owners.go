package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsOrphan(ownerRefs []metav1.OwnerReference, ownerKind string) bool {
	return len(ownerRefs) == 0 || ownerRefs[0].Kind != ownerKind
}
