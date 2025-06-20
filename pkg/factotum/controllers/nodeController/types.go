package nodecontroller

import (
	"context"
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/rjbrown57/factotum/api/v1alpha1"
	fc "github.com/rjbrown57/factotum/pkg/factotum"
	fcHandlers "github.com/rjbrown57/factotum/pkg/factotum/handlers"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const controllerName = "nodeController"

type NodeController struct {
	K8sClient   *kubernetes.Clientset
	Watcher     watch.Interface
	MsgChan     chan NcMsg
	Wg          *sync.WaitGroup
	NodeConfigs map[string]*v1alpha1.NodeConfig // a cache for updates triggered by the watcher
	NcMu        *sync.Mutex                     //NodeConfig Mutex
	NodeCache   *Cache
	Handlers    []fc.Handler
}

func NewNodeController(k8sClient *kubernetes.Clientset) (*NodeController, error) {

	// Initialize the NodeController with a Kubernetes client
	// and an empty map of NodeLabels
	nc := &NodeController{
		K8sClient:   k8sClient,
		NodeConfigs: make(map[string]*v1alpha1.NodeConfig),
		MsgChan:     make(chan NcMsg),
		Wg:          &sync.WaitGroup{},
		NodeCache: &Cache{
			ObjMap: make(map[string]*v1.Node),
			Mu:     &sync.Mutex{},
		},
		NcMu: &sync.Mutex{},
		Handlers: []fc.Handler{
			&fcHandlers.MetaDataHandler{},
			&TaintHandler{},
		},
	}

	log.Info("Initializing", "Controller", controllerName)

	// Set up a watch on the nodes in the cluster
	watcher, err := k8sClient.CoreV1().Nodes().Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error(err, "Error setting watch on nodes")
		return nil, err
	}

	nc.Watcher = watcher

	// Start watching for node events
	debugLog.Info("Starting to WatchNodes routine")
	go nc.Watch(watcher.ResultChan())

	// Start the NodeApplier that will apply labels to nodes
	debugLog.Info("Starting ApplyLabels routine")
	go nc.Proccessor()

	return nc, nil
}
