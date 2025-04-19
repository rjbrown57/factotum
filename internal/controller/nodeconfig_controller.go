/*
Copyright 2025 rjbrown57.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"slices"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	factotumiov1alpha1 "github.com/rjbrown57/factotum/api/v1alpha1"
	"github.com/rjbrown57/factotum/pkg/k8s"
	nc "github.com/rjbrown57/factotum/pkg/nodeController"
)

// NodeConfigReconciler reconciles a NodeConfig object
type NodeConfigReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	K8sClient *kubernetes.Clientset
	// We keep a copy of all existing node labels in the cluster
	NodeConfigs map[string]*factotumiov1alpha1.NodeConfig
	Nc          *nc.NodeController
}

// +kubebuilder:rbac:groups=factotum.io,resources=nodeconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=factotum.io,resources=nodeconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=factotum.io,resources=nodeconfigs/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *NodeConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	controllerLog := log.FromContext(ctx)
	DebugLog := controllerLog.V(1)

	controllerLog.Info("Reconciling NodeConfig", "name", req.NamespacedName.String())

	// Fetch the NodeConfig instance
	nodeConfig := &factotumiov1alpha1.NodeConfig{}

	if err := r.Get(ctx, req.NamespacedName, nodeConfig); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the NodeConfig is being deleted
	if !nodeConfig.DeletionTimestamp.IsZero() {
		// Handle deletion logic
		DebugLog.Info("NodeConfig is being deleted", "name", req.NamespacedName.Name)

		// Cleanup the NodeConfig instance
		// This will remove all labels, annotations, and taints from the NodeConfig
		// When passed to NodeUpdate, it will remove all labels, annotations, and taints from the node
		nodeConfig.Cleanup()
		r.Nc.NcMu.Lock()
		r.NodeConfigs[req.NamespacedName.String()] = nodeConfig
		r.Nc.NcMu.Unlock()

		// Cleanup up the NodeConfig instance
		r.Nc.Notify(nc.NcMsg{
			Header: "Cleanup",
			Node:   nil,
			Config: nodeConfig.DeepCopy(),
		})

		// Wait for the NodeController to finish processing
		r.Nc.Wg.Wait()

		r.Nc.NcMu.Lock()
		delete(r.NodeConfigs, req.NamespacedName.String())
		r.Nc.NcMu.Unlock()

		// Remove the finalizer from the NodeConfig
		nodeConfig.RemoveFinalizer()
		if err := r.Update(ctx, nodeConfig); err != nil {
			controllerLog.Error(err, "Unable to update NodeConfig with finalizer")
			return ctrl.Result{
				Requeue: true,
			}, err
		}

		controllerLog.Info("Removed finalizer from NodeConfig", "name", req.NamespacedName.String())
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	// This will prevent the CR from being deleted until we remove the finalizer
	if !slices.Contains(nodeConfig.GetFinalizers(), factotumiov1alpha1.FinalizerName) {
		nodeConfig.SetFinalizers(append(nodeConfig.GetFinalizers(), factotumiov1alpha1.FinalizerName))
		if err := r.Update(ctx, nodeConfig); err != nil {
			controllerLog.Error(err, "Unable to update NodeConfig with finalizer")
			return ctrl.Result{
				Requeue: true,
			}, err
		}
		controllerLog.Info("Added finalizer to NodeConfig", "name", req.NamespacedName.Name)
	}

	// The NodeConfig instance is being created or updated
	// We need to update the NodeConfig instance in the map
	DebugLog.Info("NodeConfig found, updating map", "name", req.NamespacedName, "labels", nodeConfig.Spec.Labels)
	r.Nc.NcMu.Lock()
	r.NodeConfigs[req.NamespacedName.String()] = nodeConfig
	r.Nc.NcMu.Unlock()

	// Send a message to the NodeController to process the config
	DebugLog.Info("Sending message to NodeController to apply configs", "NodeConfigs", len(r.NodeConfigs))
	r.Nc.Notify(nc.NcMsg{
		Header: "Reconciler",
		Node:   nil,
		Config: nodeConfig,
	})

	// Wait for the NodeController to finish processing
	r.Nc.Wg.Wait()

	controllerLog.Info("Reconciling NodeConfig complete", "name", req.NamespacedName.String())

	// Update the status of the NodeConfig
	r.Nc.NcMu.Lock()
	nodeConfig.UpdateStatus()
	r.Nc.NcMu.Unlock()

	return ctrl.Result{}, r.Status().Update(ctx, nodeConfig)
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&factotumiov1alpha1.NodeConfig{}).
		Complete(r)

	if err != nil {
		return err
	}

	r.NodeConfigs = make(map[string]*factotumiov1alpha1.NodeConfig)

	r.K8sClient = k8s.NewK8sClient()
	r.Nc, err = nc.NewNodeController(r.K8sClient)
	if err != nil {
		return err
	}

	r.Nc.NodeConfigs = r.NodeConfigs

	return nil
}
