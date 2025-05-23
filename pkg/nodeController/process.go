package nodecontroller

import (
	factotum "github.com/rjbrown57/factotum/api/v1alpha1"
	"github.com/rjbrown57/factotum/pkg/k8s"

	v1 "k8s.io/api/core/v1"
)

func (nc *NodeController) Update(node *v1.Node, NodeConfig *factotum.NodeConfig) error {

	var update bool = false
	var err error = nil

	for _, h := range nc.Handlers {
		// Call the handler functions
		if h.Update(node, NodeConfig) {
			update = true
		}
	}

	if update {
		_, err = k8s.UpdateNode(nc.K8sClient, node)
		if err != nil {
			DebugLog.Error(err, "Error updating node", "node", node.Name)
		} else {
			DebugLog.Info("Updated node", "node", node.Name)
		}
	}

	return err
}

// Proccessor will apply the changes to the nodes
// It will be called when the NodeController receives a message on the receive only MsgChan channel
func (nc *NodeController) Proccessor() error {

	for msg := range nc.MsgChan {

		// If msg node is nil, we apply to all nodes, This indicates the msg is from the reconciler and we can pass the config
		switch {
		case msg.Node == nil:
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

func (nc *NodeController) GetMatchingNodeConfigs(node *v1.Node) []*factotum.NodeConfig {
	var matchingConfigs []*factotum.NodeConfig

	nc.NcMu.Lock()

	for _, NodeConfig := range nc.NodeConfigs {
		if NodeConfig.Match(node) {
			matchingConfigs = append(matchingConfigs, NodeConfig)
		}
	}

	nc.NcMu.Unlock()
	return matchingConfigs
}

func (nc *NodeController) GetMatchingNodes(NodeConfig *factotum.NodeConfig) []*v1.Node {
	var matchingNodes []*v1.Node

	for _, node := range nc.NodeCache.NodeMap {
		if NodeConfig.Match(node) {
			matchingNodes = append(matchingNodes, node)
		}
	}

	return matchingNodes
}
