package handlers

import (
	"github.com/rjbrown57/factotum/pkg/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rjbrown57/factotum/pkg/factotum"
)

type MetaDataHandler struct{}

// MetaDataHandler will update the metadata of the object
// based on the annotations and labels defined in the FactotumConfig
// It will return true if any of the metadata was updated
func (m *MetaDataHandler) Update(Object any, FactotumConfig factotum.Config) bool {

	v1Object, ok := Object.(v1.Object)

	if !ok {
		return false
	}

	update := false

	if k8s.ProcessMetaDataMap(v1Object.GetAnnotations(), FactotumConfig.GetAnnotationSet()) {
		update = true
	}

	if k8s.ProcessMetaDataMap(v1Object.GetLabels(), FactotumConfig.GetLabelSet()) {
		update = true
	}

	return update
}
