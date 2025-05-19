package nodecontroller

import (
	"github.com/rjbrown57/factotum/api/v1alpha1"
	"github.com/rjbrown57/factotum/pkg/k8s"

	v1 "k8s.io/api/core/v1"
)

func (nc *NodeController) Update(node *v1.Node, NodeConfig *v1alpha1.NodeConfig) error {

	var err error = nil

	newNode := node.DeepCopy()

	for _, h := range nc.Handlers {
		// Call the handler functions
		_ = h.Update(newNode, NodeConfig)
	}

	_, err = k8s.StrategicMerge(nc.K8sClient, node, newNode)
	if err != nil {
		DebugLog.Error(err, "Error updating node", "node", node.Name)
	} else {
		DebugLog.Info("Updated node", "node", node.Name)
	}

	return err
}

// Proccessor will apply the changes to the nodes
// It will be called when the NodeController receives a message on the receive only MsgChan channel
func (nc *NodeController) Proccessor() error {

	for msg := range nc.MsgChan {

		// If msg node is nil, this indicates a Config Event
		// Process all matching nodes
		switch {
		case msg.Node == nil:

			if msg.Config.DetectChange() {
				log.Info("Selector has changed, processing all nodes", "config", msg.Config.Spec.Selector)
				// Make a deep copy of the config to avoid modifying the original
				c := msg.Config.DeepCopy()
				// Call cleanup to remove any config from the no longer selected nodes
				c.Cleanup()

				for _, node := range nc.GetNodeDiffSet(c.Status.AppliedSelector, c.Spec.Selector) {
					log.Info("Processing node", "node", node.Name)
					if err := nc.Update(node, c); err != nil {
						log.Error(err, "Error processing node", "node", node.Name)
					}
				}
			}

			for _, node := range nc.GetMatchingNodes(msg.Config) {
				log.Info("Processing node", "node", node.Name)
				if err := nc.Update(node, msg.Config); err != nil {
					log.Error(err, "Error processing node", "node", node.Name)
				}
			}

		// Update to a specific node
		// If msg node is not nil, we apply to the specific node, This indicates the msg is from the watcher so we need to use our cache
		case msg.Node != nil:
			node := msg.Node
			// If the Node Has Configs that match we will process the node
			for _, NodeConfig := range nc.GetMatchingNodeConfigs(node) {
				log.Info("Processing node", "node", node.Name)
				if err := nc.Update(node, NodeConfig); err != nil {
					log.Error(err, "Error processing node", "node", node.Name)
				}
			}
		}

		// Notify the WaitGroup that we are done processing
		nc.Wg.Done()
	}

	return nil
}

func (nc *NodeController) GetMatchingNodeConfigs(node *v1.Node) []*v1alpha1.NodeConfig {
	var matchingConfigs []*v1alpha1.NodeConfig

	nc.NcMu.Lock()

	for _, NodeConfig := range nc.NodeConfigs {
		if NodeConfig.Match(node) {
			matchingConfigs = append(matchingConfigs, NodeConfig)
		}
	}

	nc.NcMu.Unlock()
	return matchingConfigs
}

// GetNodeSet returns the nodes matching the selectors
func (nc *NodeController) GetMatchingNodes(NodeConfig *v1alpha1.NodeConfig) []*v1.Node {
	var matchingNodes []*v1.Node

	for _, node := range nc.NodeCache.ObjMap {
		if matchNode(node, NodeConfig.Spec.Selector) {
			matchingNodes = append(matchingNodes, node)
		}
	}

	return matchingNodes
}

// GetNodeSet returns the nodes No longer matching the selector, but matched the previous selector
func (nc *NodeController) GetNodeDiffSet(PreviousSelector, CurrentSelector v1alpha1.NodeSelector) []*v1.Node {
	var matchingNodes []*v1.Node

	for _, node := range nc.NodeCache.ObjMap {
		if matchNode(node, PreviousSelector) && !matchNode(node, CurrentSelector) {
			matchingNodes = append(matchingNodes, node)
		}
	}

	return matchingNodes
}

// matchNode checks if a node matches the given selector
func matchNode(node *v1.Node, selector v1alpha1.NodeSelector) bool {
	if selector.NodeSelector == nil {
		return true
	}

	for key, value := range selector.NodeSelector {
		if nodeValue, exists := node.Labels[key]; !exists || nodeValue != value {
			return false
		}
	}

	return true
}
