package nodecontroller

import (
	"sync"

	v1 "k8s.io/api/core/v1"
)

type Cache struct {
	ObjMap map[string]*v1.Node
	Mu     *sync.Mutex
}

func (Cache *Cache) GetNode(name string) (*v1.Node, bool) {
	Cache.Mu.Lock()
	defer Cache.Mu.Unlock()

	node, ok := Cache.ObjMap[name]
	if !ok {
		return nil, false
	}

	return node, true
}

func (Cache *Cache) SetNode(name string, node *v1.Node) {
	Cache.Mu.Lock()
	defer Cache.Mu.Unlock()

	Cache.ObjMap[name] = node.DeepCopy()
}

func (Cache *Cache) DeleteNode(name string) {
	Cache.Mu.Lock()
	defer Cache.Mu.Unlock()

	delete(Cache.ObjMap, name)
}
