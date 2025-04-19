package handlers

import (
	"reflect"
	"slices"

	factotum "github.com/rjbrown57/factotum/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type TaintHandler struct{}

func (t *TaintHandler) Update(node *v1.Node, NodeConfig *factotum.NodeConfig) bool {

	update := false

	nodeTaintMap := SliceToMap(node.Spec.Taints)

	// needs to be replaced with get Taint set
	for _, taint := range NodeConfig.GetTaintSet() {
		switch currentTaint, exists := nodeTaintMap[taint.Key]; {
		// Taint is missing in node, add it
		case !exists:
			node.Spec.Taints = append(node.Spec.Taints, taint)
			update = true
		// Taint is wrong in node, update it
		case !reflect.DeepEqual(currentTaint, taint):
			if taint.Effect == "" {
				if index := FindTaintIndex(taint.Key, node.Spec.Taints); index != -1 {
					node.Spec.Taints = slices.Delete(node.Spec.Taints, index, index+1)
				}
				update = true
			} else {
				if index := FindTaintIndex(taint.Key, node.Spec.Taints); index != -1 {
					node.Spec.Taints[index] = taint
					update = true
				}
			}

		}
	}

	return update
}

func SliceToMap(taints []v1.Taint) map[string]v1.Taint {
	taintMap := make(map[string]v1.Taint)
	for _, taint := range taints {
		taintMap[taint.Key] = taint
	}
	return taintMap
}

func FindTaintIndex(key string, taints []v1.Taint) int {
	for i, nodeTaint := range taints {
		if nodeTaint.Key == key {
			return i
		}
	}
	// If the taint is not found, return -1
	return -1
}
