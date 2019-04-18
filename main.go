package main

import (
	"flag"
	"time"

	ctrl "github.com/christopherhein/node-operator/controller"
	clientset "github.com/christopherhein/node-operator/generated/clientset/versioned"
	informers "github.com/christopherhein/node-operator/generated/informers/externalversions"
	"github.com/christopherhein/node-operator/signals"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var masterURL string
var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig")
	flag.StringVar(&masterURL, "master", "", "The url of the API server.")
}

func main() {
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("error building config error=%s", err.Error())
	}

	nodeClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("error creating node client error=%s", err.Error())
	}

	clientInformerFactory := informers.NewSharedInformerFactory(nodeClient, time.Minute*30)

	controller := ctrl.New(nodeClient, clientInformerFactory.Node().V1alpha1().AuthorizedKeys())

	clientInformerFactory.Start(stopCh)

	if err = controller.Run(1, stopCh); err != nil {
		klog.Fatalf("error running controller error=%s", err.Error())
	}
}
