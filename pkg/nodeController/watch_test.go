package nodecontroller

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCompareNodes(t *testing.T) {
	tests := []struct {
		name     string
		node1    *v1.Node
		node2    *v1.Node
		expected bool
	}{
		{
			name: "Nodes are equal",
			node1: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"key1": "value1"},
					Labels:      map[string]string{"label1": "value1"},
				},
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			node2: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"key1": "value1"},
					Labels:      map[string]string{"label1": "value1"},
				},
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			expected: true,
		},
		{
			name: "Nodes have different annotations",
			node1: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"key1": "value1"},
				},
			},
			node2: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"key1": "value2"},
				},
			},
			expected: false,
		},
		{
			name: "Nodes have different labels",
			node1: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"label1": "value1"},
				},
			},
			node2: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"label1": "value2"},
				},
			},
			expected: false,
		},
		{
			name: "Nodes have different taints",
			node1: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value1", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			node2: &v1.Node{
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						{Key: "key1", Value: "value2", Effect: v1.TaintEffectNoSchedule},
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareNodes(tt.node1, tt.node2)
			if result != tt.expected {
				t.Errorf("CompareNodes() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
