package namespacecontroller

import (
	"sync"

	v1 "k8s.io/api/core/v1"
)

type Cache struct {
	ObjMap map[string]*v1.Namespace
	Mu     *sync.Mutex
}

func (Cache *Cache) Get(name string) (*v1.Namespace, bool) {
	Cache.Mu.Lock()
	defer Cache.Mu.Unlock()

	obj, ok := Cache.ObjMap[name]
	if !ok {
		return nil, false
	}

	return obj, true
}

func (Cache *Cache) Set(name string, obj *v1.Namespace) {
	Cache.Mu.Lock()
	defer Cache.Mu.Unlock()

	Cache.ObjMap[name] = obj.DeepCopy()

	debugLog.Info("NS Cache Set NS", "ns", name)

}

func (Cache *Cache) Delete(name string) {
	Cache.Mu.Lock()
	defer Cache.Mu.Unlock()

	delete(Cache.ObjMap, name)
	debugLog.Info("Ns Cache Delete NS", "ns", name)

}
