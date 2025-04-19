package handlers

import (
	factotum "github.com/rjbrown57/factotum/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type Handler interface {
	Update(node *v1.Node, NodeConfig *factotum.NodeConfig) bool
}
