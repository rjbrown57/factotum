package nodecontroller

import (
	"sync"

	v1 "k8s.io/api/core/v1"
)

type NodeCache struct {
	NodeMap map[string]*v1.Node
	nodeMu  *sync.Mutex
}

func (NodeCache *NodeCache) GetNode(name string) (*v1.Node, bool) {
	NodeCache.nodeMu.Lock()
	defer NodeCache.nodeMu.Unlock()

	node, ok := NodeCache.NodeMap[name]
	if !ok {
		return nil, false
	}

	return node, true
}

func (NodeCache *NodeCache) SetNode(name string, node *v1.Node) {
	NodeCache.nodeMu.Lock()
	defer NodeCache.nodeMu.Unlock()

	NodeCache.NodeMap[name] = node.DeepCopy()
}

func (NodeCache *NodeCache) DeleteNode(name string) {
	NodeCache.nodeMu.Lock()
	defer NodeCache.nodeMu.Unlock()

	delete(NodeCache.NodeMap, name)
}
