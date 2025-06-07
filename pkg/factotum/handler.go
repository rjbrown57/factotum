package factotum

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Handers are called with an object and a FactotumConfig
// and are expected to update the object based on the config
// and return true if any of the metadata was updated
type Handler interface {
	Update(object v1.Object, FactotumConfig Config) v1.Object
	GetName() string
}
