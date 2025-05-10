package v1alpha1

import (
	"testing"

	"github.com/rjbrown57/factotum/pkg/factotum/config"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRemoveFinalizer(t *testing.T) {
	nc := &NodeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Finalizers: []string{FinalizerName, "other-finalizer"},
		},
	}

	nc.RemoveFinalizer()

	if len(nc.Finalizers) != 1 || nc.Finalizers[0] != "other-finalizer" {
		t.Errorf("expected finalizer to be removed, got %v", nc.Finalizers)
	}
}

func TestGetLabelSet(t *testing.T) {
	nc := &NodeConfig{
		Spec: NodeConfigSpec{
			CommonSpec: config.CommonSpec{
				Labels: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
		Status: NodeConfigStatus{
			CommonStatus: config.CommonStatus{
				AppliedLabels: map[string]string{
					"key2": "value2",
					"key3": "value3",
				}},
		},
	}

	labelSet := nc.GetLabelSet()

	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "",
	}

	for k, v := range expected {
		if labelSet[k] != v {
			t.Errorf("expected %s=%s, got %s=%s", k, v, k, labelSet[k])
		}
	}
}

func TestMatch(t *testing.T) {

	nc := &NodeConfig{
		Spec: NodeConfigSpec{
			Selector: NodeSelector{
				NodeSelector: map[string]string{
					"key1": "value1.*",
					"key2": "value2",
				},
			},
		},
	}

	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"key1": "value123",
				"key2": "value2",
			},
		},
	}

	if !nc.Match(node) {
		t.Errorf("expected node to match NodeConfig")
	}

	node.Labels["key2"] = "wrongValue"
	if nc.Match(node) {
		t.Errorf("expected node not to match NodeConfig")
	}
}

func TestErrorStatus(t *testing.T) {
	nc := &NodeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-nodeconfig",
		},
		Spec: NodeConfigSpec{
			CommonSpec: config.CommonSpec{
				Labels:      map[string]string{"key1": "value1"},
				Annotations: map[string]string{"annotation1": "value1"},
			},
		},
	}

	nc.ErrorStatus()

	if len(nc.Status.Conditions) != 1 {
		t.Errorf("expected 1 condition, got %d", len(nc.Status.Conditions))
	}

	condition := nc.Status.Conditions[0]
	if condition.Type != "Applied" || condition.Status != metav1.ConditionFalse {
		t.Errorf("unexpected condition: %+v", condition)
	}
}

func TestUpdateStatus(t *testing.T) {
	nc := &NodeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-nodeconfig",
		},
		Spec: NodeConfigSpec{
			CommonSpec: config.CommonSpec{
				Labels:      map[string]string{"key1": "value1"},
				Annotations: map[string]string{"annotation1": "value1"},
			},
		},
	}

	nc.UpdateStatus()

	if len(nc.Status.Conditions) != 1 {
		t.Errorf("expected 1 condition, got %d", len(nc.Status.Conditions))
	}

	condition := nc.Status.Conditions[0]
	if condition.Type != "Applied" || condition.Status != metav1.ConditionTrue {
		t.Errorf("unexpected condition: %+v", condition)
	}
}

func TestCleanup(t *testing.T) {
	nc := &NodeConfig{
		Spec: NodeConfigSpec{
			CommonSpec: config.CommonSpec{
				Labels: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				Annotations: map[string]string{
					"annotation1": "value1",
					"annotation2": "value2",
				},
			},
			Taints: []corev1.Taint{
				{
					Key:    "key1",
					Effect: corev1.TaintEffectNoSchedule,
					Value:  "value1",
				},
			},
		},
	}

	nc.Cleanup()

	if len(nc.Spec.Labels) != 0 && len(nc.Annotations) != 0 && len(nc.Spec.Taints) != 0 {
		t.Errorf("expected all configs to be empty, got %v", nc)
	}
}

func TestGetTaintSet(t *testing.T) {
	nc := &NodeConfig{
		Spec: NodeConfigSpec{
			Taints: []corev1.Taint{
				{
					Key:    "key1",
					Effect: corev1.TaintEffectNoSchedule,
					Value:  "value1",
				},
			},
		},
		Status: NodeConfigStatus{
			AppliedTaints: []corev1.Taint{
				{
					Key:    "key2",
					Effect: corev1.TaintEffectNoSchedule,
					Value:  "value2",
				},
			},
		},
	}

	expected := nc.GetTaintSet()

	if len(expected) != 2 {
		t.Errorf("expected 2 taints, got %d", len(expected))
	}
}
