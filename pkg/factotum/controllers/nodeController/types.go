package nodecontroller

import (
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/rjbrown57/factotum/api/v1alpha1"
	fc "github.com/rjbrown57/factotum/pkg/factotum"
	fcHandlers "github.com/rjbrown57/factotum/pkg/factotum/handlers"
)

const controllerName = "nodeController"

type NodeController struct {
	K8sClient   *kubernetes.Clientset
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

	// Start watching for node events
	debugLog.Info("Starting to WatchNodes routine")

	// Start the NodeApplier that will apply labels to nodes
	debugLog.Info("Starting ApplyLabels routine")
	go nc.Proccessor()

	return nc, nil
}
