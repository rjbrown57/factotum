package namespacecontroller

import (
	"github.com/rjbrown57/factotum/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type Msg struct {
	Header    string
	Namespace *v1.Namespace
	Config    *v1alpha1.NamespaceConfig
}

func (c *NamespaceController) Notify(msg Msg) {
	log.Info("Notifying NamespaceController", "source", msg.Header)
	c.Wg.Add(1)
	c.MsgChan <- msg
}
