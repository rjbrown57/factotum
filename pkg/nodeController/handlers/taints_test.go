package handlers

import (
	"testing"

	factotum "github.com/rjbrown57/factotum/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func TestTaintHandler_Update(t *testing.T) {

	tests := []struct {
		name           string
		node           *v1.Node
		nodeConfig     *factotum.NodeConfig
		expectedTaints []v1.Taint
		expectedUpdate bool
	}{
		{
			name: "Add missing taint",
			node: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{},
				},
			},
			nodeConfig: &factotum.NodeConfig{
				Spec: factotum.NodeConfigSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			expectedTaints: []v1.Taint{
				{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
			},
			expectedUpdate: true,
		},
		{
			name: "Remove taint when effect is empty",
			node: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			nodeConfig: &factotum.NodeConfig{
				Spec: factotum.NodeConfigSpec{
					Taints: []v1.Taint{
						{Key: "key1", Effect: ""},
					},
				},
			},
			expectedTaints: []v1.Taint{},
			expectedUpdate: true,
		},
		{
			name: "Update taint value",
			node: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "oldValue", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			nodeConfig: &factotum.NodeConfig{
				Spec: factotum.NodeConfigSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "newValue", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			expectedTaints: []v1.Taint{
				{Key: "key1", Value: "newValue", Effect: v1.TaintEffectNoSchedule},
			},
			expectedUpdate: true,
		},
		{
			name: "No changes when taints match",
			node: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			nodeConfig: &factotum.NodeConfig{
				Spec: factotum.NodeConfigSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			expectedTaints: []v1.Taint{
				{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
			},
			expectedUpdate: false,
		},
		{
			name: "Add multiple missing taints",
			node: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{},
				},
			},
			nodeConfig: &factotum.NodeConfig{
				Spec: factotum.NodeConfigSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
						{Key: "key2", Value: "value2", Effect: v1.TaintEffectPreferNoSchedule},
					},
				},
			},
			expectedTaints: []v1.Taint{
				{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
				{Key: "key2", Value: "value2", Effect: v1.TaintEffectPreferNoSchedule},
			},
			expectedUpdate: true,
		},
		{
			name: "Remove taint not in config",
			node: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
						{Key: "key2", Value: "value2", Effect: v1.TaintEffectPreferNoSchedule},
					},
				},
			},
			nodeConfig: &factotum.NodeConfig{
				Spec: factotum.NodeConfigSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
				Status: factotum.NodeConfigStatus{
					AppliedTaints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
						{Key: "key2", Value: "value2", Effect: v1.TaintEffectPreferNoSchedule},
					},
				},
			},
			expectedTaints: []v1.Taint{
				{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
			},
			expectedUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TaintHandler{}
			updated := handler.Update(tt.node, tt.nodeConfig)

			if updated != tt.expectedUpdate {
				t.Errorf("expected update %v, got %v", tt.expectedUpdate, updated)
			}

			if len(tt.node.Spec.Taints) != len(tt.expectedTaints) {
				t.Errorf("expected taints length %d, got %d, %v", len(tt.expectedTaints), len(tt.node.Spec.Taints), tt.node.Spec.Taints)
			}

			for i, expectedTaint := range tt.expectedTaints {
				if tt.node.Spec.Taints[i] != expectedTaint {
					t.Errorf("expected taint %v, got %v", expectedTaint, tt.node.Spec.Taints[i])
				}
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
