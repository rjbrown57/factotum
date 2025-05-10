package config

import (
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const FinalizerName string = "factotum.io/factotum"

// +k8s:deepcopy-gen=true
type CommonSpec struct {
	// Annotations to Apply to Selected Objects
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels to Apply to Selected Objects
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

func RemoveFinalizer(m *metav1.ObjectMeta) {
	for i, finalizer := range m.GetFinalizers() {
		if finalizer == FinalizerName {
			m.SetFinalizers(slices.Delete(m.Finalizers, i, i+1))
			break
		}
	}
}

func ProcessMap(desiredMap, statusMap map[string]string) map[string]string {
	for key := range statusMap {
		if _, exists := desiredMap[key]; !exists {
			desiredMap[key] = ""
		}
	}

	return desiredMap
}

// CommonStatus defines the observed state of a Config
// +k8s:deepcopy-gen=true
type CommonStatus struct {
	// Conditions is an array of conditions that describe the status of the object
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Labels applied to the objects
	AppliedLabels map[string]string `json:"appliedLabels,omitempty"`
	// Annotations applied to the objects
	AppliedAnnotations map[string]string `json:"appliedAnnotations,omitempty"`
}
