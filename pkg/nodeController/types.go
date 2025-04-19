package nodecontroller

import (
	"context"
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	factotum "github.com/rjbrown57/factotum/api/v1alpha1"
	nodeHandlers "github.com/rjbrown57/factotum/pkg/nodeController/handlers"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var log = ctrl.Log.WithName("nc")
var DebugLog = log.V(1)

type NodeController struct {
	K8sClient   *kubernetes.Clientset
	Watcher     watch.Interface
	MsgChan     chan NcMsg
	Wg          *sync.WaitGroup
	NodeConfigs map[string]*factotum.NodeConfig // a cache for updates triggered by the watcher
	NcMu        *sync.Mutex                     //NodeConfig Mutex
	NodeCache   *NodeCache
	Handlers    []nodeHandlers.Handler
}

func NewNodeController(k8sClient *kubernetes.Clientset) (*NodeController, error) {

	log.Info("Initializing NodeController")
	// Initialize the NodeController with a Kubernetes client
	// and an empty map of NodeLabels
	nc := &NodeController{
		K8sClient:   k8sClient,
		NodeConfigs: make(map[string]*factotum.NodeConfig),
		MsgChan:     make(chan NcMsg),
		Wg:          &sync.WaitGroup{},
		NodeCache: &NodeCache{
			NodeMap: make(map[string]*v1.Node),
			nodeMu:  &sync.Mutex{},
		},
		NcMu: &sync.Mutex{},
		Handlers: []nodeHandlers.Handler{
			&nodeHandlers.MetaDataHandler{},
			&nodeHandlers.TaintHandler{},
		},
	}

	// Set up a watch on the nodes in the cluster
	watcher, err := k8sClient.CoreV1().Nodes().Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error(err, "Error setting watch on nodes")
		return nil, err
	}

	nc.Watcher = watcher

	// Start watching for node events
	DebugLog.Info("Starting to WatchNodes routine")
	go nc.Watch(watcher.ResultChan())

	// Start the NodeApplier that will apply labels to nodes
	DebugLog.Info("Starting ApplyLabels routine")
	go nc.Proccessor()

	return nc, nil
}
