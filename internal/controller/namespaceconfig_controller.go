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

	"github.com/rjbrown57/factotum/api/v1alpha1"
	"github.com/rjbrown57/factotum/pkg/factotum/config"
	controller "github.com/rjbrown57/factotum/pkg/factotum/controllers/namespaceController"
	"github.com/rjbrown57/factotum/pkg/k8s"
)

// NamespaceConfigReconciler reconciles a NamespaceConfig object
type NamespaceConfigReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	k8sClient       *kubernetes.Clientset
	NamspaceConfigs map[string]*v1alpha1.NamespaceConfig
	Controller      *controller.NamespaceController
}

// +kubebuilder:rbac:groups=factotum.io,resources=namespaceconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=factotum.io,resources=namespaceconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=factotum.io,resources=namespaceconfigs/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;update;create;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NamespaceConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *NamespaceConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	controllerLog.Info("Reconciling NamespaceConfig", "name", req.NamespacedName.String())

	// Fetch the NamespaceConfig instance
	fConfig := &v1alpha1.NamespaceConfig{}

	if err := r.Get(ctx, req.NamespacedName, fConfig); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the NamespaceConfig is being deleted
	if !fConfig.DeletionTimestamp.IsZero() {
		// Handle deletion logic
		debugLog.Info("NamespaceConfig is being deleted", "name", req.NamespacedName.Name)

		// Cleanup the NamespaceConfig instance
		// This will remove all labels, annotations, and taints from the NamespaceConfig
		// When passed to NodeUpdate, it will remove all labels, annotations, and taints from the node
		fConfig.Cleanup()
		r.Controller.Mu.Lock()
		r.NamspaceConfigs[req.NamespacedName.String()] = fConfig
		r.Controller.Mu.Unlock()

		// Cleanup up the NamespaceConfig instance
		r.Controller.Notify(controller.Msg{
			Header:    "Cleanup",
			Namespace: nil,
			Config:    fConfig.DeepCopy(),
		})

		// Wait for the NodeController to finish processing
		r.Controller.Wg.Wait()

		r.Controller.Mu.Lock()
		delete(r.NamspaceConfigs, req.NamespacedName.String())
		r.Controller.Mu.Unlock()

		// Remove the finalizer from the NamespaceConfig
		fConfig.RemoveFinalizer()
		if err := r.Update(ctx, fConfig); err != nil {
			controllerLog.Error(err, "Unable to update NamespaceConfig with finalizer")
			return ctrl.Result{
				Requeue: true,
			}, err
		}

		controllerLog.Info("Removed finalizer from NamespaceConfig", "name", req.NamespacedName.String())
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	// This will prevent the CR from being deleted until we remove the finalizer
	if !slices.Contains(fConfig.GetFinalizers(), config.FinalizerName) {
		fConfig.SetFinalizers(append(fConfig.GetFinalizers(), config.FinalizerName))
		if err := r.Update(ctx, fConfig); err != nil {
			controllerLog.Error(err, "Unable to update NamespaceConfig with finalizer")
			return ctrl.Result{
				Requeue: true,
			}, err
		}
		controllerLog.Info("Added finalizer to NamespaceConfig", "name", req.NamespacedName.Name)
	}

	// The NamespaceConfig instance is being created or updated
	// We need to update the NamespaceConfig instance in the map
	debugLog.Info("NamespaceConfig found, updating map", "name", req.NamespacedName, "labels", fConfig.Spec.Labels)
	r.Controller.Mu.Lock()
	r.NamspaceConfigs[req.NamespacedName.String()] = fConfig
	r.Controller.Mu.Unlock()

	// Send a message to the NodeController to process the config
	debugLog.Info("Sending message to NodeController to apply configs", "NamespaceConfigs", len(r.NamspaceConfigs))
	r.Controller.Notify(controller.Msg{
		Header:    "Reconciler",
		Namespace: nil,
		Config:    fConfig,
	})

	// Wait for the NodeController to finish processing
	r.Controller.Wg.Wait()

	controllerLog.Info("Reconciling NamespaceConfig complete", "name", req.NamespacedName.String())

	// Update the status of the NamespaceConfig
	r.Controller.Mu.Lock()
	fConfig.UpdateStatus()
	r.Controller.Mu.Unlock()

	return ctrl.Result{}, r.Status().Update(ctx, fConfig)
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NamespaceConfig{}).
		Complete(r)

	if err != nil {
		return err
	}

	r.NamspaceConfigs = make(map[string]*v1alpha1.NamespaceConfig)

	r.k8sClient = k8s.NewK8sClient()
	r.Controller, err = controller.NewNamespaceController(r.k8sClient, r.NamspaceConfigs)

	return err
}
