// Package main for the node operator
package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"

	nodeioclient "github.com/christopherhein/node-operator/pkg/client/clientset/versioned/typed/node/v1"
	nodeioAuthorizedKey "github.com/christopherhein/node-operator/pkg/operator/authorizedkey"
	opkit "github.com/christopherhein/operator-kit"
	"k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig file")

	log.Info("Getting kubernetes context")
	context, nodeClientset, err := createContext(*kubeconfig)
	if err != nil {
		log.Fatalf("failed to create context. %+v\n", err)
	}

	// Create and wait for CRD resources
	log.Info("Registering the authorized key resource")
	resources := []opkit.CustomResource{
		nodeioAuthorizedKey.AuthorizedKeyResource,
	}
	err = opkit.CreateCustomResources(*context, resources)
	if err != nil {
		log.Fatalf("failed to create custom resource. %+v\n", err)
	}

	// create signals to stop watching the resources
	signalChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start watching the authorized key resource
	log.Info("Watching the authorized key resource")
	controller := nodeioAuthorizedKey.NewController(context, nodeClientset)
	controller.StartWatch(v1.NamespaceAll, stopChan)

	for {
		select {
		case <-signalChan:
			log.Info("shutdown signal received, exiting...")
			close(stopChan)
			return
		}
	}
}

func getClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func createContext(kubeconfig string) (*opkit.Context, nodeioclient.NodeV1Interface, error) {
	config, err := getClientConfig(kubeconfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get k8s config. %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get k8s client. %+v", err)
	}

	apiExtClientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create k8s API extension clientset. %+v", err)
	}

	nodeClientset, err := nodeioclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create node clientset. %+v", err)
	}

	context := &opkit.Context{
		Clientset:             clientset,
		APIExtensionClientset: apiExtClientset,
		Interval:              500 * time.Millisecond,
		Timeout:               60 * time.Second,
	}
	return context, nodeClientset, nil
}
