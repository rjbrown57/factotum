package nodecontroller

import (
	"reflect"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Watch Nodes will keep our node list up to date
// On change it will notify the NodeController to update labels
// MsgChan is a send only channel that will be used to notify the NodeController
// ch is a receive only channel that will be used to receive events from the watch
func (nc *NodeController) Watch(ch <-chan watch.Event) error {

	for event := range ch {
		debugLog.Info("Node Watcher", "event", event.Type)
		switch event.Type {
		case watch.Added, watch.Modified:

			node, ok := event.Object.(*v1.Node)
			if !ok {
				log.Error(nil, "Error casting event object to Node")
				continue
			}

			newNode, exists := nc.NodeCache.GetNode(node.Name)
			if !exists {
				// If the node doesn't exist in the cache, add it
				nc.NodeCache.SetNode(node.Name, node)
				continue
			}

			if !CompareNodes(node, newNode) {
				nc.Notify(NcMsg{
					Header: "Watcher",
					Node:   node,
				})
			}

			// we always update the cache even if the node changed a field we don't care about
			nc.NodeCache.SetNode(node.Name, node)

		case watch.Deleted:
			node, ok := event.Object.(*v1.Node)
			if !ok {
				log.Error(nil, "Error casting event object to Node")
				continue
			}
			// Remove the node from the node list
			nc.NodeCache.DeleteNode(node.Name)
		}
	}

	return nil
}

// CompareNodes compares two nodes and returns true if they are equal
func CompareNodes(node1, node2 *v1.Node) bool {

	if !reflect.DeepEqual(node1.Annotations, node2.Annotations) {
		traceLog.Info("Node Annotations differ", "node1", node1.Name, "node2", node2.Name)
		return false
	}

	if !reflect.DeepEqual(node1.Labels, node2.Labels) {
		traceLog.Info("Node Labels differ", "node1", node1.Name, "node2", node2.Name)
		return false
	}

	if !reflect.DeepEqual(node1.Spec.Taints, node2.Spec.Taints) {
		traceLog.Info("Node Taints differ", "node1", node1.Name, "node2", node2.Name)
		return false
	}

	return true
}
