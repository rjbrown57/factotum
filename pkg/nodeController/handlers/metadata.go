package handlers

import (
	factotum "github.com/rjbrown57/factotum/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type MetaDataHandler struct{}

func (m *MetaDataHandler) Update(node *v1.Node, NodeConfig *factotum.NodeConfig) bool {

	update := false

	if ProcessMetaDataMap(node.Annotations, NodeConfig.GetAnnotationSet()) {
		update = true
	}

	if ProcessMetaDataMap(node.Labels, NodeConfig.GetLabelSet()) {
		update = true
	}

	return update
}

// This should probably be moved to k8s package since it's common across
func ProcessMetaDataMap(currentMap, desiredMap map[string]string) bool {
	update := false

	for key, value := range desiredMap {
		switch currentValue, exists := currentMap[key]; {
		// Label is empty in config, remove it from node
		case value == "":
			// Label is empty, remove it
			delete(currentMap, key)
			update = true
		// Label is missing in node, add it
		case !exists:
			currentMap[key] = value
			update = true
		// Label is wrong in node, update it
		case currentValue != value:
			currentMap[key] = value
			update = true
		}
	}

	return update
}
