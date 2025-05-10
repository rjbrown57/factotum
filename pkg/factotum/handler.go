package factotum

// Handers are called with an object and a FactotumConfig
// and are expected to update the object based on the config
// and return true if any of the metadata was updated
type Handler interface {
	Update(object any, FactotumConfig Config) bool
}
