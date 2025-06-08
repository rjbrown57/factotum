package namespacecontroller

import (
	"reflect"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Watch Nodes will keep our obj list up to date
// On change it will notify the NamespaceController to update labels
// MsgChan is a send only channel that will be used to notify the NamespaceController
// ch is a receive only channel that will be used to receive events from the watch
func (c *NamespaceController) Watch(ch <-chan watch.Event) error {

	for event := range ch {
		debugLog.Info("NS Cache Watcher", "event", event.Type)
		switch event.Type {
		case watch.Added, watch.Modified:

			obj, ok := event.Object.(*v1.Namespace)
			if !ok {
				log.Error(nil, "Error casting event object to NS")
				continue
			}

			newNode, exists := c.Cache.Get(obj.Name)
			if !exists {
				// If the obj doesn't exist in the cache, add it
				c.Cache.Set(obj.Name, obj)
				continue
			}

			if !Compare(obj, newNode) {
				c.Notify(Msg{
					Header:    "Watcher",
					Namespace: obj,
				})
			}

			// we always update the cache even if the obj changed a field we don't care about
			c.Cache.Set(obj.Name, obj)

		case watch.Deleted:
			obj, ok := event.Object.(*v1.Namespace)
			if !ok {
				log.Error(nil, "Error casting event object to Node")
				continue
			}
			// Remove the obj from the obj list
			c.Cache.Delete(obj.Name)
		}
	}

	return nil
}

// Compare compares two objs and returns true if they are equal based on fields we care about
func Compare(obj1, obj2 *v1.Namespace) bool {

	if !reflect.DeepEqual(obj1.Annotations, obj2.Annotations) {
		traceLog.Info("Annotations not equal", "obj1", obj1.Name, "obj2", obj2.Name)
		return false
	}

	if !reflect.DeepEqual(obj1.Labels, obj2.Labels) {
		traceLog.Info("Labels not equal", "obj1", obj1.Name, "obj2", obj2.Name)
		return false
	}

	return true
}
