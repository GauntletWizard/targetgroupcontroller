package main

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"net/http"
)

// NewClient is a basic wrapper for the boilerplate of creating a kubernetes client.
// It's basically just the example client.
func NewK8sClient() (client *kubernetes.Clientset, namespace string) {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &overrides, os.Stdin)

	conf, err := clientConfig.ClientConfig()
	if err != nil {
		log.Panic("Failed to initialize clientConfig, ", err)
	}
	ns, _, err := clientConfig.Namespace()
	if err != nil {
		log.Panic("Failed to initialize clientConfig Namespace, ", err)
	}
	client = kubernetes.NewForConfigOrDie(conf)
	log.Println("Client started")
	return client, ns
}
