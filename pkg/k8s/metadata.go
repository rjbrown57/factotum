package k8s

import "fmt"

// ProcessMetaDataMap processes the currentMap and desiredMap for metadata updates.
func ProcessMetaDataMap(currentMap, desiredMap map[string]string) map[string]string {

	// Check if the desiredMap is nil, if so, initialize it
	if desiredMap == nil {
		desiredMap = make(map[string]string)
	}

	// Check if the currentMap is nil, if so, initialize it
	if currentMap == nil {
		return desiredMap
	}

	for key, value := range desiredMap {
		switch currentValue, exists := currentMap[key]; {
		// Label is empty in config, remove it from node
		case value == "":
			// Label is empty, remove it
			delete(currentMap, key)
			fmt.Printf("Removing %s from obj\n", key)
		// Label is missing in node, add it
		case !exists:
			currentMap[key] = value
			fmt.Printf("Adding %s to obj\n", key)
		// Label is wrong in node, update it
		case currentValue != value:
			currentMap[key] = value
			fmt.Printf("Updating %s to %s in obj\n", key, value)
		}
	}

	return currentMap
}
