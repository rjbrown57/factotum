/*
Copyright 2025 rjbrown57.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	"regexp"

	"slices"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// https://book.kubebuilder.io/reference/markers/crd-validation
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodeConfigSpec defines the desired state of NodeConfig
type NodeConfigSpec struct {
	// Annotations to Apply to Selected Nodes, If no selector is provided, all nodes will be selected
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels to Apply to Selected Nodes, If no selector is provided, all nodes will be selected
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Taints to Apply to Selected Nodes, If no selector is provided, all nodes will be selected
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`

	// NodeSelector is a map of node labels to select nodes
	// +optional
	Selector NodeSelector `json:"selector"`
}

// NodeConfigStatus defines the observed state of NodeConfig
type NodeConfigStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Labels applied to the nodes
	AppliedLabels map[string]string `json:"appliedLabels,omitempty"`
	// Annotations applied to the nodes
	AppliedAnnotations map[string]string `json:"appliedAnnotations,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AppliedTaints []corev1.Taint `json:"appliedTaints,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// NodeConfig is the Schema for the nodeconfigs API
type NodeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeConfigSpec   `json:"spec,omitempty"`
	Status NodeConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NodeConfigList contains a list of NodeConfig
type NodeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeConfig{}, &NodeConfigList{})
}

type NodeSelector struct {
	// NodeSelector is a map of node labels to select nodes
	// Selector can be provided a plain string or a regex.
	// If no selector is provided, all nodes will be selected
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

const FinalizerName string = "factotum.io/factotum"

func (nc *NodeConfig) RemoveFinalizer() {
	for i, finalizer := range nc.GetFinalizers() {
		if finalizer == FinalizerName {
			nc.SetFinalizers(slices.Delete(nc.Finalizers, i, i+1))
			break
		}
	}
}

// Cleanup removes all labels, annotations, and taints from the NodeConfig
// When passed to NodeUpdate, it will remove all labels, annotations, and taints from the node
func (nc *NodeConfig) Cleanup() {
	nc.Spec.Labels = make(map[string]string)
	nc.Spec.Annotations = make(map[string]string)
	nc.Spec.Taints = make([]corev1.Taint, 0)
}

// GetLabelSet compares the labels in the NodeConfig with the labels in the appliedLabels status
// and returns a map of labels that need to be applied to the nodes.
func (nc *NodeConfig) GetLabelSet() map[string]string {

	// If the labels are nil, create a new map
	if nc.Spec.Labels == nil {
		nc.Spec.Labels = make(map[string]string)
	}

	labelSet := nc.Spec.Labels

	// If the label is in the appliedLabels status, but not in the spec, remove it from the labelSet
	for key := range nc.Status.AppliedLabels {
		if _, exists := nc.Spec.Labels[key]; !exists {
			labelSet[key] = ""
		}
	}

	return labelSet
}

// These two are basically identical, so we should replace them with a single function
func (nc *NodeConfig) GetAnnotationSet() map[string]string {

	// If the annotations are nil, create a new map
	if nc.Spec.Annotations == nil {
		nc.Spec.Annotations = make(map[string]string)
	}

	AnnotationSet := nc.Spec.Annotations

	for key := range nc.Status.AppliedAnnotations {
		if _, exists := nc.Spec.Annotations[key]; !exists {
			AnnotationSet[key] = ""
		}
	}

	return AnnotationSet
}

func (nc *NodeConfig) ErrorStatus() {
	nc.Status.AppliedLabels = nc.Spec.Labels
	nc.Status.AppliedAnnotations = nc.Spec.Annotations
	nc.Status.AppliedTaints = nc.Spec.Taints
	nc.Status.Conditions = []metav1.Condition{
		{
			Type:               "Applied",
			Status:             metav1.ConditionFalse,
			Reason:             "NodeConfigError",
			Message:            fmt.Sprintf("%s MalFormed NodeConfig", fmt.Sprintf("%s/%s", nc.Namespace, nc.Name)),
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: nc.Generation,
		},
	}
}

func (nc *NodeConfig) UpdateStatus() {
	nc.Status.AppliedLabels = nc.Spec.Labels
	nc.Status.AppliedAnnotations = nc.Spec.Annotations
	nc.Status.AppliedTaints = nc.Spec.Taints
	nc.Status.Conditions = []metav1.Condition{
		{
			Type:               "Applied",
			Status:             metav1.ConditionTrue,
			Reason:             "NodeConfigReady",
			Message:            fmt.Sprintf("%s Applied", fmt.Sprintf("%s/%s", nc.Namespace, nc.Name)),
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: nc.Generation,
		},
	}
}

// Match checks if the node matches all selectors in the NodeConfig
func (nc *NodeConfig) Match(node *corev1.Node) bool {
	if nc.Spec.Selector.NodeSelector == nil {
		return true
	}

	for SelectorKey, SelectorValue := range nc.Spec.Selector.NodeSelector {

		//  All Selector Labels must match
		if _, exists := node.Labels[SelectorKey]; !exists {
			return false
		}

		if match, err := regexp.MatchString(SelectorValue, node.Labels[SelectorKey]); err != nil || !match {
			// If the regex does not match, return false
			return false

		}
	}

	return true
}

// WIP will come back to this

func (nc *NodeConfig) GetTaintSet() []corev1.Taint {
	// If the taints are nil, create a new slice
	if nc.Spec.Taints == nil {
		nc.Spec.Taints = make([]corev1.Taint, 0)
	}

	taintSet := nc.Spec.Taints

	// The possible scenarios are
	// 1. The taint is in AppliedTaints and in the spec, so do nothing
	// 2. The taint is in the spec, but not in the appliedTaints, so add it to the taintSet handled on 216
	// 3. The taint is in the appliedTaints, but not in the spec, mark for removal

	// If the taint is in the appliedTaints status, but not in the spec, remove it from the taintSet
	for _, taint := range nc.Status.AppliedTaints {
		if _, exists := nc.FindTaint(taint.Key); !exists {
			taint.Effect = ""
			taintSet = append(taintSet, taint)
		}
	}

	return taintSet
}

// Match checks if the node matches all selectors in the NodeConfig
func (nc *NodeConfig) FindTaint(key string) (v1.Taint, bool) {
	for _, taint := range nc.Spec.Taints {
		if taint.Key == key {
			return taint, true
		}
	}
	return v1.Taint{}, false
}
