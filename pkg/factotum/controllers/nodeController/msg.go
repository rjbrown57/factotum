package nodecontroller

import (
	"github.com/rjbrown57/factotum/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type NcMsg struct {
	Header string
	Node   *v1.Node
	Config *v1alpha1.NodeConfig
}

func (nc *NodeController) Notify(msg NcMsg) {
	log.Info("Notifying NodeController", "source", msg.Header)
	nc.Wg.Add(1)
	nc.MsgChan <- msg
}
