package nodecontroller

import (
	"testing"

	"github.com/rjbrown57/factotum/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func TestTaintHandler_Update(t *testing.T) {
	handler := TaintHandler{}

	tests := []struct {
		name           string
		nodeConfig     *v1alpha1.NodeConfig
		initialObject  *v1.Node
		expectedObject *v1.Node
	}{
		{
			name: "Add new taint",
			nodeConfig: &v1alpha1.NodeConfig{
				Spec: v1alpha1.NodeConfigSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			initialObject: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{},
				},
			},
			expectedObject: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
		},
		{
			name: "Update existing taint",
			nodeConfig: &v1alpha1.NodeConfig{
				Spec: v1alpha1.NodeConfigSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "newValue", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			initialObject: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			expectedObject: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "newValue", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.Update(tt.initialObject, tt.nodeConfig)
			if result == nil {
				t.Errorf("expected non-nil result")
				return
			}
			node, ok := result.(*v1.Node)
			if !ok {
				t.Errorf("expected result to be of type *v1.Node")
				return
			}
			if len(node.Spec.Taints) != len(tt.expectedObject.Spec.Taints) {
				t.Errorf("expected %d taints, got %d", len(tt.expectedObject.Spec.Taints), len(node.Spec.Taints))
			}
		})
	}

}

func TestFindTaintIndex(t *testing.T) {
	tests := []struct {
		name      string
		taints    []v1.Taint
		key       string
		expected  int
		expectErr bool
	}{
		{
			name: "Taint exists",
			taints: []v1.Taint{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			key:      "key2",
			expected: 1,
		},
		{
			name: "Taint does not exist",
			taints: []v1.Taint{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			key:       "key3",
			expected:  -1,
			expectErr: true,
		},
		{
			name:      "Empty taints list",
			taints:    []v1.Taint{},
			key:       "key1",
			expected:  -1,
			expectErr: true,
		},
		{
			name: "1 taints",
			taints: []v1.Taint{
				{Key: "key1", Value: "value1"},
			},
			key:       "key1",
			expected:  0,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := FindTaintIndex(tt.key, tt.taints)
			if index != tt.expected {
				t.Errorf("expected index %d, got %d", tt.expected, index)
			}
			if (index == -1) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, index == -1)
			}
		})
	}
}
func TestSliceToMap(t *testing.T) {
	tests := []struct {
		name     string
		taints   []v1.Taint
		expected map[string]v1.Taint
	}{
		{
			name: "Non-empty taints list",
			taints: []v1.Taint{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			expected: map[string]v1.Taint{
				"key1": {Key: "key1", Value: "value1"},
				"key2": {Key: "key2", Value: "value2"},
			},
		},
		{
			name:     "Empty taints list",
			taints:   []v1.Taint{},
			expected: map[string]v1.Taint{},
		},
		{
			name: "Duplicate keys in taints list",
			taints: []v1.Taint{
				{Key: "key1", Value: "value1"},
				{Key: "key1", Value: "value2"},
			},
			expected: map[string]v1.Taint{
				"key1": {Key: "key1", Value: "value2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SliceToMap(tt.taints)
			if len(result) != len(tt.expected) {
				t.Errorf("expected map length %d, got %d", len(tt.expected), len(result))
			}
			for key, expectedTaint := range tt.expected {
				if result[key] != expectedTaint {
					t.Errorf("for key %s, expected taint %v, got %v", key, expectedTaint, result[key])
				}
			}
		})
	}
}
