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
func (m *MetaDataHandler) Update(Object v1.Object, FactotumConfig factotum.Config) v1.Object {

	_ = k8s.ProcessMetaDataMap(Object.GetAnnotations(), FactotumConfig.GetAnnotationSet())
	_ = k8s.ProcessMetaDataMap(Object.GetLabels(), FactotumConfig.GetLabelSet())

	return Object
}
