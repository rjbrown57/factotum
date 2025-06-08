package namespacecontroller

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

const controllerName = "namespacecontroller"

type NamespaceController struct {
	K8sClient        *kubernetes.Clientset
	Watcher          watch.Interface
	MsgChan          chan Msg
	Wg               *sync.WaitGroup
	NamespaceConfigs map[string]*v1alpha1.NamespaceConfig // a cache for updates triggered by the watcher
	Mu               *sync.Mutex                          //NamespaceConfig Mutex
	Cache            *Cache
	Handlers         []fc.Handler
}

func NewNamespaceController(k8sClient *kubernetes.Clientset, SharedCache map[string]*v1alpha1.NamespaceConfig) (*NamespaceController, error) {

	log.Info("Initializing", "Controller", controllerName)

	// Initialize the NamespaceController with a Kubernetes client
	// and an empty map of NodeLabels
	c := &NamespaceController{
		K8sClient:        k8sClient,
		NamespaceConfigs: make(map[string]*v1alpha1.NamespaceConfig),
		MsgChan:          make(chan Msg),
		Wg:               &sync.WaitGroup{},
		Cache: &Cache{
			ObjMap: make(map[string]*v1.Namespace),
			Mu:     &sync.Mutex{},
		},
		Mu: &sync.Mutex{},
		Handlers: []fc.Handler{
			&fcHandlers.MetaDataHandler{},
		},
	}

	// Set up a watch on the objs in the cluster
	watcher, err := k8sClient.CoreV1().Namespaces().Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error(err, "Error setting watch on objs")
		return nil, err
	}

	c.Watcher = watcher

	// Start watching for obj events
	debugLog.Info("Starting to Watch routine")
	go c.Watch(watcher.ResultChan())

	// Start the NodeApplier that will apply labels to objs
	debugLog.Info("Starting ApplyLabels routine")
	go c.Proccessor()

	return c, nil
}
