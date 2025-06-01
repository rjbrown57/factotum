package k8s

import (
	"reflect"
	"testing"
)

func TestProcessMetaDataMap(t *testing.T) {
	tests := []struct {
		name        string
		currentMap  map[string]string
		desiredMap  map[string]string
		expectedMap map[string]string
		updated     bool
	}{
		{
			name:        "Add new key-value pair",
			currentMap:  map[string]string{"key1": "value1"},
			desiredMap:  map[string]string{"key1": "value1", "key2": "value2"},
			expectedMap: map[string]string{"key1": "value1", "key2": "value2"},
			updated:     true,
		},
		{
			name:        "Update existing key-value pair",
			currentMap:  map[string]string{"key1": "value1"},
			desiredMap:  map[string]string{"key1": "newValue"},
			expectedMap: map[string]string{"key1": "newValue"},
			updated:     true,
		},
		{
			name:        "Remove key with empty value",
			currentMap:  map[string]string{"key1": "value1", "key2": "value2"},
			desiredMap:  map[string]string{"key1": "value1", "key2": ""},
			expectedMap: map[string]string{"key1": "value1"},
			updated:     true,
		},
		{
			name:        "No changes needed",
			currentMap:  map[string]string{"key1": "value1"},
			desiredMap:  map[string]string{"key1": "value1"},
			expectedMap: map[string]string{"key1": "value1"},
			updated:     false,
		},
		{
			name:        "Remove all keys",
			currentMap:  map[string]string{"key1": "value1", "key2": "value2"},
			desiredMap:  map[string]string{"key1": "", "key2": ""},
			expectedMap: map[string]string{},
			updated:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentMap := tt.currentMap
			updated := ProcessMetaDataMap(currentMap, tt.desiredMap)

			if updated != tt.updated {
				t.Errorf("expected updated to be %v, got %v", tt.updated, updated)
			}

			if !reflect.DeepEqual(currentMap, tt.expectedMap) {
				t.Errorf("expected map to be %v, got %v", tt.expectedMap, currentMap)
			}
		})
	}
}
