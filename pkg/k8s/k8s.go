package k8s

import (
	"context"
	"errors"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"encoding/json"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

func GetConfig() *rest.Config {

	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	// If we are not in a cluster, we will try to use the kubeconfig file
	if errors.Is(err, rest.ErrNotInCluster) {

		// This should get smarter to detect alternative kubeconfig locations
		// break it into it's own function
		kubeconfig := fmt.Sprintf("%s/.kube/config", homedir.HomeDir())

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error creating kubeconfig: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Error getting in-cluster config: %v", err)
	}

	return config
}

func NewK8sClient() *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(GetConfig())
	if err != nil {
		log.Fatalf("Error creating kubernetes clientset: %v", err)
	}
	return clientset
}

func GetClientset(c *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		log.Fatalf("Error creating kubernetes clientset: %v", err)
	}
	return clientset, nil
}

func GetNodes(c *kubernetes.Clientset) (*v1.NodeList, error) {
	return c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

func StrategicMerge(c *kubernetes.Clientset, original metav1.Object, modified metav1.Object) (metav1.Object, error) {
	// Use the appropriate client to apply the patch based on the object's type
	// This example assumes the object is a Namespace, but you should handle other types as needed

	originalJSON, err := json.Marshal(original)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal original object: %w", err)
	}

	modifiedJSON, err := json.Marshal(modified)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified object: %w", err)
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(originalJSON, modifiedJSON, original)

	switch obj := original.(type) {
	case *v1.Namespace:
		if err != nil {
			return nil, fmt.Errorf("failed to create strategic merge patch: %w", err)
		}
		return c.CoreV1().Namespaces().Patch(context.TODO(), obj.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	case *v1.Node:
		if err != nil {
			return nil, fmt.Errorf("failed to create strategic merge patch: %w", err)
		}
		return c.CoreV1().Nodes().Patch(context.TODO(), obj.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	default:
		return nil, fmt.Errorf("unsupported object type")
	}
}
