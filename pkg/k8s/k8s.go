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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func UpdateNode(c *kubernetes.Clientset, node *v1.Node) (*v1.Node, error) {
	return c.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
}
