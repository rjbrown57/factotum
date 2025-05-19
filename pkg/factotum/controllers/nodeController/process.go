package nodecontroller

import (
	"reflect"

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

			// If AppliedSelector is different we need to clean the diff between applied and current
			if !reflect.DeepEqual(msg.Config.Status.AppliedSelector, msg.Config.Spec.Selector) {
				log.Info("Selector Change detected", "nodeConfig", msg.Config.Name)

				// This will work but create a bunch of cycles

				// Calculate the diff set
				// Make a copy and update the selector to our "known" selector
				// Same logic needs to be done for the labelSet / annotationSet / taints
				diffSet := msg.Config.DeepCopy()
				diffSet.Spec.Selector = msg.Config.Status.AppliedSelector

				//oldNodes := nc.GetMatchingNodes(diffSet)
				//newNodes := nc.GetMatchingNodes(msg.Config)

				// compare them and find the diff
				nodes := []*v1.Node{}

				for _, node := range nodes {
					if err := nc.Update(node, msg.Config); err != nil {
					}

					// now that we have the diff range

					// if there are any presetn in the applied set that are not in the current set
					// we need to clean them
				}

				for _, node := range nc.GetMatchingNodes(msg.Config) {
					log.Info("Processing node", "node", node.Name)
					if err := nc.Update(node, msg.Config); err != nil {
						log.Error(err, "Error processing node", "node", node.Name)
					}
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

func (nc *NodeController) GetMatchingNodes(NodeConfig *v1alpha1.NodeConfig) []*v1.Node {
	var matchingNodes []*v1.Node

	for _, node := range nc.NodeCache.ObjMap {
		if NodeConfig.Match(node) {
			matchingNodes = append(matchingNodes, node)
		}
	}

	return matchingNodes
}
