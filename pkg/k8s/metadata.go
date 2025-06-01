package k8s

// This should probably be moved to k8s package since it's common across
func ProcessMetaDataMap(currentMap, desiredMap map[string]string) bool {
	update := false

	// Check if the currentMap is nil, if so, initialize it
	if currentMap == nil {
		currentMap = make(map[string]string)
	}

	// Check if the desiredMap is nil, if so, initialize it
	if desiredMap == nil {
		desiredMap = make(map[string]string)
	}

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
