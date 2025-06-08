package handlers

import (
	"github.com/rjbrown57/factotum/pkg/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rjbrown57/factotum/pkg/factotum"
)

type MetaDataHandler struct{}

// MetaDataHandler will update the metadata of the object
// based on the annotations and labels defined in the FactotumConfig
func (m *MetaDataHandler) Update(Object v1.Object, FactotumConfig factotum.Config) v1.Object {

	Object.SetAnnotations(k8s.ProcessMetaDataMap(Object.GetAnnotations(), FactotumConfig.GetAnnotationSet()))
	Object.SetLabels(k8s.ProcessMetaDataMap(Object.GetLabels(), FactotumConfig.GetLabelSet()))

	return Object
}

func (m *MetaDataHandler) GetName() string {
	return "MetaDataHandler"
}
