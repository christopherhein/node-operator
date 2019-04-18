package controller

import (
	"fmt"
	"time"

	nodev1alpha1 "github.com/christopherhein/node-operator/apis/node/v1alpha1"
	"github.com/christopherhein/node-operator/authorizedkey"
	clientset "github.com/christopherhein/node-operator/generated/clientset/versioned"
	nodescheme "github.com/christopherhein/node-operator/generated/clientset/versioned/scheme"
	informers "github.com/christopherhein/node-operator/generated/informers/externalversions/node/v1alpha1"
	listers "github.com/christopherhein/node-operator/generated/listers/node/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

const controllerName = "node-controller"

type Controller struct {
	nodeClientset clientset.Interface
	authKeyLister listers.AuthorizedKeyLister
	authKeySynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
}

func New(nodeClientset clientset.Interface, authKeyInformer informers.AuthorizedKeyInformer) *Controller {
	utilruntime.Must(nodescheme.AddToScheme(scheme.Scheme))

	controller := &Controller{
		nodeClientset: nodeClientset,
		authKeyLister: authKeyInformer.Lister(),
		authKeySynced: authKeyInformer.Informer().HasSynced,
		workqueue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "AuthorizedKeys"),
	}

	klog.Info("setting up event handlers")

	authKeyInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueAuthKey,
		UpdateFunc: func(old, new interface{}) {
			newKey := new.(*nodev1alpha1.AuthorizedKey)
			oldKey := old.(*nodev1alpha1.AuthorizedKey)
			if newKey.ResourceVersion == oldKey.ResourceVersion {
				return
			}
			controller.enqueueAuthKey(new)
		},
		DeleteFunc: controller.deleteAuthKey,
	})

	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("starting node controller")

	klog.Info("waiting for caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.authKeySynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("started workers")
	<-stopCh
	klog.Info("shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing %s=%s, requeuing", key, err.Error())
		}

		c.workqueue.Forget(obj)
		klog.Infof("successfully synced %s", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key=%s", key))
		return nil
	}

	authorizedKeyFile := authorizedkey.File{
		UID: key,
	}

	authKey, err := c.authKeyLister.Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("authorized key %s in workqueue no longer exists", key))
			authorizedKeyFile.Sync(true)
			return nil
		}
		return err
	}

	authorizedKeyFile = authorizedkey.File{
		UID: authKey.Name,
		Key: authKey.Data.Key,
	}

	err = authorizedKeyFile.Sync(false)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) enqueueAuthKey(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

func (c *Controller) deleteAuthKey(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}
