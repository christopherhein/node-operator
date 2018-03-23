package authorizedkey

import (
	log "github.com/sirupsen/logrus"
	"reflect"

	opkit "github.com/christopherhein/operator-kit"
	// apiv1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	nodeio "github.com/christopherhein/node-operator/pkg/apis/node.io"
	nodeioV1 "github.com/christopherhein/node-operator/pkg/apis/node.io/v1"
	nodeioclient "github.com/christopherhein/node-operator/pkg/client/clientset/versioned/typed/node/v1"
)

// AuthorizedKeyResource is the Authorized Key CRD definition
var AuthorizedKeyResource = opkit.CustomResource{
	Name:       "authorizedkey",
	Plural:     "authorizedkeys",
	Group:      nodeio.GroupName,
	Version:    nodeio.Version,
	Scope:      apiextensionsv1beta1.NamespaceScoped,
	Kind:       reflect.TypeOf(nodeioV1.AuthorizedKey{}).Name(),
	ShortNames: []string{"authkey", "authorized-key", "authorized-keys"},
}

// Controller represents a controller object for object store custom resources
type Controller struct {
	context         *opkit.Context
	nodeioClientset nodeioclient.NodeV1Interface
}

// NewController create controller for watching object store custom resources created
func NewController(context *opkit.Context, nodeioClientset nodeioclient.NodeV1Interface) *Controller {
	return &Controller{
		context:         context,
		nodeioClientset: nodeioClientset,
	}
}

// StartWatch watches for instances of Object Store custom resources and acts on them
func (c *Controller) StartWatch(namespace string, stopCh chan struct{}) error {
	resourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	}
	restClient := c.nodeioClientset.RESTClient()
	watcher := opkit.NewWatcher(AuthorizedKeyResource, namespace, resourceHandlers, restClient)
	go watcher.Watch(&nodeioV1.AuthorizedKey{}, stopCh)
	return nil
}

func (c *Controller) onAdd(obj interface{}) {
	object := obj.(*nodeioV1.AuthorizedKey).DeepCopy()

	authorizedKeyFile := File{
		UID: string(object.ObjectMeta.UID),
		Key: object.Data.Key,
	}
	err := authorizedKeyFile.WriteKey()
	if err != nil {
		log.Errorf("Could not add key '%s', errored with %+v", object.Name, err)
	}

	log.Infof("Added authorized key '%s'", object.Name)
}

func (c *Controller) onUpdate(oldObj, newObj interface{}) {
	oldObject := oldObj.(*nodeioV1.AuthorizedKey).DeepCopy()
	newObject := newObj.(*nodeioV1.AuthorizedKey).DeepCopy()

	authorizedKeyFile := File{
		UID: string(oldObject.ObjectMeta.UID),
		Key: newObject.Data.Key,
	}
	err := authorizedKeyFile.UpdateKey()
	if err != nil {
		log.Errorf("Could not update key '%s', errored with %+v", newObject.Name, err)
	}

	log.Infof("Updated authorized key '%s'", newObject.Name)
}

func (c *Controller) onDelete(obj interface{}) {
	object := obj.(*nodeioV1.AuthorizedKey).DeepCopy()

	authorizedKeyFile := File{
		UID: string(object.ObjectMeta.UID),
		Key: object.Data.Key,
	}
	err := authorizedKeyFile.DeleteKey()
	if err != nil {
		log.Errorf("Could not delete key '%s', errored with %+v", object.Name, err)
	}

	log.Infof("Deleted authorized key '%s'", object.Name)
}
