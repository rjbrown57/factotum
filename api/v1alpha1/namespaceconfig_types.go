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

	"github.com/rjbrown57/factotum/pkg/factotum/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NamespaceConfigSpec defines the desired state of NamespaceConfig
type NamespaceConfigSpec struct {
	config.CommonSpec `json:",inline"`
	Selector          NamespaceSelector   `json:"selector,omitempty"`
	Namespaces        []FactotumNamespace `json:"namespaces,omitempty"`
}

// NamespaceConfigStatus defines the observed state of NamespaceConfig
type NamespaceConfigStatus struct {
	config.CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// NamespaceConfig is the Schema for the namespaceconfigs API
type NamespaceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceConfigSpec   `json:"spec,omitempty"`
	Status NamespaceConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NamespaceConfigList contains a list of NamespaceConfig
type NamespaceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NamespaceConfig{}, &NamespaceConfigList{})
}

type NamespaceSelector struct {
	NamespaceSelector map[string]string `json:"namespaceSelector,omitempty"`
}

// FactotumNamespaces are namespaces that are managed by Factotum
// They are used to create namespaces with specific labels and annotations
type FactotumNamespace struct {
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

func (nc *NamespaceConfig) RemoveFinalizer() {
	config.RemoveFinalizer(&nc.ObjectMeta)
}

// Cleanup removes all labels, annotations, and taints from the NamespaceConfig
// When passed to NodeUpdate, it will remove all labels, annotations, and taints from the node
func (c *NamespaceConfig) Cleanup() {
	c.Spec.Labels = make(map[string]string)
	c.Spec.Annotations = make(map[string]string)
}

// GetLabelSet compares the labels in the NamespaceConfig with the labels in the appliedLabels status
// and returns a map of labels that need to be applied to the nodes.
func (c *NamespaceConfig) GetLabelSet() map[string]string {

	// If the labels are nil, create a new map
	if c.Spec.Labels == nil {
		c.Spec.Labels = make(map[string]string)
	}

	return config.ProcessMap(c.Spec.Labels, c.Status.AppliedLabels)
}

// These two are basically identical, so we should replace them with a single function
func (nc *NamespaceConfig) GetAnnotationSet() map[string]string {

	// If the annotations are nil, create a new map
	if nc.Spec.Annotations == nil {
		nc.Spec.Annotations = make(map[string]string)
	}

	return config.ProcessMap(nc.Spec.Annotations, nc.Status.AppliedAnnotations)
}

func (c *NamespaceConfig) ErrorStatus() {
	c.Status.AppliedLabels = c.Spec.Labels
	c.Status.AppliedAnnotations = c.Spec.Annotations
	c.Status.Conditions = []metav1.Condition{
		{
			Type:               "Applied",
			Status:             metav1.ConditionFalse,
			Reason:             "NamespaceConfigError",
			Message:            fmt.Sprintf("%s MalFormed NamespaceConfig", fmt.Sprintf("%s/%s", c.Namespace, c.Name)),
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: c.Generation,
		},
	}
}

func (c *NamespaceConfig) UpdateStatus() {

	// Clean will remove all empty labels and annotations from the NamespaceConfig
	c.Spec.Clean()

	c.Status.AppliedLabels = c.Spec.Labels
	c.Status.AppliedAnnotations = c.Spec.Annotations
	c.Status.Conditions = []metav1.Condition{
		{
			Type:               "Applied",
			Status:             metav1.ConditionTrue,
			Reason:             "NamespaceConfigReady",
			Message:            fmt.Sprintf("%s Applied", fmt.Sprintf("%s/%s", c.Namespace, c.Name)),
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: c.Generation,
		},
	}
}

// Match checks if the node matches all selectors in the NamespaceConfig
func (nc *NamespaceConfig) Match(obj *corev1.Namespace) bool {

	if nc.Spec.Selector.NamespaceSelector == nil {
		return true
	}

	for SelectorKey, SelectorValue := range nc.Spec.Selector.NamespaceSelector {

		//  All Selector Labels must match
		if _, exists := obj.Labels[SelectorKey]; !exists {
			return false
		}

		if match, err := regexp.MatchString(SelectorValue, obj.Labels[SelectorKey]); err != nil || !match {
			// If the regex does not match, return false
			return false

		}
	}

	return true
}
