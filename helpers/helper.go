package helpers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// isKubeNamespace checks if the given namespace is a Kubernetes-owned namespace.
func IsKubeNamespace(ns string) bool {
	return ns == metav1.NamespacePublic || ns == metav1.NamespaceSystem
}
