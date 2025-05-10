package factotum

import (
	"k8s.io/apimachinery/pkg/watch"
)

type FactotumController interface {
	Watch(ch <-chan watch.Event) error
	Processor() error
	Cleanup() error
	AddFinalizer() error
	UpdateCache()
}
