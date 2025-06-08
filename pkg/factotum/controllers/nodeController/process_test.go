package nodecontroller

import (
	"testing"

	"github.com/rjbrown57/factotum/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeNode(name string, labels map[string]string) *v1.Node {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
}

func makeNodeConfig(selector map[string]string) *v1alpha1.NodeConfig {
	return &v1alpha1.NodeConfig{
		Spec: v1alpha1.NodeConfigSpec{
			Selector: v1alpha1.NodeSelector{
				NodeSelector: selector,
			},
		},
	}
}

func TestMatchNode(t *testing.T) {
	node := makeNode("node1", map[string]string{"foo": "bar", "baz": "qux"})

	tests := []struct {
		name     string
		selector v1alpha1.NodeSelector
		want     bool
	}{
		{
			name:     "single label match",
			selector: v1alpha1.NodeSelector{NodeSelector: map[string]string{"foo": "bar"}},
			want:     true,
		},
		{
			name:     "multiple labels all match",
			selector: v1alpha1.NodeSelector{NodeSelector: map[string]string{"foo": "bar", "baz": "qux"}},
			want:     true,
		},
		{
			name:     "empty selector matches all",
			selector: v1alpha1.NodeSelector{NodeSelector: nil},
			want:     true,
		},
		{
			name:     "regex selector matches",
			selector: v1alpha1.NodeSelector{NodeSelector: map[string]string{"foo": "^ba.*"}},
			want:     true,
		},
		{
			name:     "no match",
			selector: v1alpha1.NodeSelector{NodeSelector: map[string]string{"zasdf": "nope"}},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchNode(node, tt.selector)
			if got != tt.want {
				t.Errorf("matchNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
