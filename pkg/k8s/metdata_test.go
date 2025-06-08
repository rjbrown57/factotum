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
	}{
		{
			name:        "Add new key-value pair",
			currentMap:  map[string]string{"key1": "value1"},
			desiredMap:  map[string]string{"key1": "value1", "key2": "value2"},
			expectedMap: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:        "Update existing key-value pair",
			currentMap:  map[string]string{"key1": "value1"},
			desiredMap:  map[string]string{"key1": "newValue"},
			expectedMap: map[string]string{"key1": "newValue"},
		},

		{
			name:        "Remove key with empty value",
			currentMap:  map[string]string{"key1": "value1", "key2": "value2"},
			desiredMap:  map[string]string{"key1": "value1", "key2": ""},
			expectedMap: map[string]string{"key1": "value1"},
		},
		{
			name:        "No changes needed",
			currentMap:  map[string]string{"key1": "value1"},
			desiredMap:  map[string]string{"key1": "value1"},
			expectedMap: map[string]string{"key1": "value1"},
		},
		{
			name:        "Remove all keys",
			currentMap:  map[string]string{"key1": "value1", "key2": "value2"},
			desiredMap:  map[string]string{"key1": "", "key2": ""},
			expectedMap: map[string]string{},
		},
		{
			name:        "Nil currentMap add new key-value pairs",
			currentMap:  nil,
			desiredMap:  map[string]string{"key1": "value1", "key2": "value2"},
			expectedMap: map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnedMap := ProcessMetaDataMap(tt.currentMap, tt.desiredMap)

			if !reflect.DeepEqual(returnedMap, tt.expectedMap) {
				t.Errorf("expected map to be %v, got %v", tt.expectedMap, returnedMap)
			}
		})
	}
}
