package namespacecontroller

import (
	"context"

	"github.com/rjbrown57/factotum/api/v1alpha1"
	"github.com/rjbrown57/factotum/pkg/k8s"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *NamespaceController) Update(namespace *v1.Namespace, NamespaceConfig *v1alpha1.NamespaceConfig) error {

	var err error = nil

	newNs := namespace.DeepCopy()

	for _, h := range c.Handlers {
		// Call the handler functions
		traceLog.Info("Calling handler", "handler", h.GetName(), "node", namespace.Name, "config", NamespaceConfig.Name)
		_ = h.Update(newNs, NamespaceConfig)
	}

	_, err = k8s.StrategicMerge(c.K8sClient, namespace, newNs)
	if err != nil {
		log.Error(err, "Error updating obj", "obj", namespace.Name)
	} else {
		log.Info("Updated obj", "obj", namespace.Name)
	}

	return err
}

// Proccessor will apply the changes to the objs
// It will be called when the NamespaceController receives a message on the receive only MsgChan channel
func (c *NamespaceController) Proccessor() error {

	for msg := range c.MsgChan {

		// If msg obj is nil, we apply to all objs, This indicates the msg is from the reconciler and we can pass the config
		switch {
		case msg.Namespace == nil:

			// First we check on if we have created all the requested namespaces

			for _, ns := range msg.Config.Spec.Namespaces {
				_, err := c.K8sClient.CoreV1().Namespaces().Get(context.TODO(), ns.Name, metav1.GetOptions{})
				switch {
				case apierrors.IsNotFound(err):
					log.Info("Namespace not found, creating", "namespace", ns.Name)
					ns := &v1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name:        ns.Name,
							Labels:      msg.Config.Spec.Labels,
							Annotations: msg.Config.Spec.Annotations,
						},
					}
					ns, err = c.K8sClient.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
					if err != nil {
						log.Error(err, "Error creating namespace", "namespace", ns.Name)
						continue
					}
					// Add the namespace to the cache
					c.Cache.ObjMap[ns.Name] = ns
				case err != nil:
					log.Error(err, "Error getting namespace", "namespace", ns.Name)
					continue
				}
			}

			for _, obj := range c.GetMatchingNamespaces(msg.Config) {
				log.Info("Processing obj", "obj", obj.Name)
				if err := c.Update(obj, msg.Config); err != nil {
					log.Error(err, "Error processing obj", "obj", obj.Name)
				}
			}
		// Update to a specific obj
		// If msg obj is not nil, we apply to the specific obj, This indicates the msg is from the watcher so we need to use our cache
		case msg.Namespace != nil:
			obj := msg.Namespace
			// If the Node Has Configs that match we will process the obj
			for _, NamespaceConfig := range c.GetMatchingNamespaceConfigs(obj) {
				log.Info("Processing obj", "obj", obj.Name)
				if err := c.Update(obj, NamespaceConfig); err != nil {
					log.Error(err, "Error processing obj", "obj", obj.Name)
				}
			}
		}

		// Notify the WaitGroup that we are done processing
		c.Wg.Done()
	}

	return nil
}

func (c *NamespaceController) GetMatchingNamespaceConfigs(obj *v1.Namespace) []*v1alpha1.NamespaceConfig {
	var matchingConfigs []*v1alpha1.NamespaceConfig

	c.Mu.Lock()

	for _, NamespaceConfig := range c.NamespaceConfigs {
		if NamespaceConfig.Match(obj) {
			matchingConfigs = append(matchingConfigs, NamespaceConfig)
		}
	}

	c.Mu.Unlock()
	return matchingConfigs
}

func (c *NamespaceController) GetMatchingNamespaces(NamespaceConfig *v1alpha1.NamespaceConfig) []*v1.Namespace {
	var matchingNamespaces []*v1.Namespace

	// Manage any created NS
	for _, obj := range NamespaceConfig.Spec.Namespaces {
		if ns, exists := c.Cache.Get(obj.Name); exists {
			matchingNamespaces = append(matchingNamespaces, ns)
		}
	}

	// Manage any selected NS
	for _, obj := range c.Cache.ObjMap {
		if NamespaceConfig.Match(obj) {
			matchingNamespaces = append(matchingNamespaces, obj)
		}
	}

	return matchingNamespaces
}
