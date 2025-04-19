package k8s

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	if config == nil {
		t.Fatalf("Expected non-nil config, got nil")
	}
}

func TestNewK8sClient(t *testing.T) {
	clientset := NewK8sClient()
	if clientset == nil {
		t.Fatalf("Expected non-nil clientset, got nil")
	}
}

func TestGetClientset(t *testing.T) {
	config := &rest.Config{
		Host: "http://localhost",
	}
	clientset, err := GetClientset(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if clientset == nil {
		t.Fatalf("Expected non-nil clientset, got nil")
	}
}

func TestGetNodes(t *testing.T) {
	clientset := NewK8sClient()
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if nodes == nil || len(nodes.Items) == 0 {
		t.Fatalf("Expected non-empty node list, got nil or empty")
	}
}
