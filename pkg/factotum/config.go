package factotum

type Config interface {
	GetAnnotationSet() map[string]string
	GetLabelSet() map[string]string
}
