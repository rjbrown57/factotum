package nodecontroller

import (
	"reflect"
	"slices"

	"github.com/rjbrown57/factotum/api/v1alpha1"
	"github.com/rjbrown57/factotum/pkg/factotum"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TaintHandler struct{}

func (t *TaintHandler) GetName() string {
	return "TaintHandler"
}

func (t *TaintHandler) Update(Object v1.Object, Config factotum.Config) v1.Object {

	node, ok := Object.(*corev1.Node)
	if !ok {
		return Object
	}

	debugLog.Info("TaintHandler Update", "node", node.Name)

	// Assert that the Config is of type NodeConfig
	// so we can access the GetTaintSet method
	NodeConfig, ok := Config.(*v1alpha1.NodeConfig)
	if !ok {
		return Object
	}

	nodeTaintMap := SliceToMap(node.Spec.Taints)

	// needs to be replaced with get Taint set
	for _, taint := range NodeConfig.GetTaintSet() {
		switch currentTaint, exists := nodeTaintMap[taint.Key]; {
		// Taint is missing in node, add it
		case !exists:
			debugLog.Info("TaintHandler Adding Taint to", "node", node.Name)
			node.Spec.Taints = append(node.Spec.Taints, taint)
		// Taint is wrong in node, update it
		case !reflect.DeepEqual(currentTaint, taint):
			debugLog.Info("TaintHandler Updating Taint on", "node", node.Name, "taint", taint.Key)
			if taint.Effect == "" {
				if index := FindTaintIndex(taint.Key, node.Spec.Taints); index != -1 {
					node.Spec.Taints = slices.Delete(node.Spec.Taints, index, index+1)
				}
			} else {
				if index := FindTaintIndex(taint.Key, node.Spec.Taints); index != -1 {
					node.Spec.Taints[index] = taint
				}
			}

		}
	}

	return Object
}

func SliceToMap(taints []corev1.Taint) map[string]corev1.Taint {
	taintMap := make(map[string]corev1.Taint)
	for _, taint := range taints {
		taintMap[taint.Key] = taint
	}
	return taintMap
}

func FindTaintIndex(key string, taints []corev1.Taint) int {
	for i, nodeTaint := range taints {
		if nodeTaint.Key == key {
			return i
		}
	}
	// If the taint is not found, return -1
	return -1
}
