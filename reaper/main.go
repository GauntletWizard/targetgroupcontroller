package main

import (
	"flag"
	"log"
	"net/http"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	client "github.com/gauntletwizard/targetgroupcontroller/k8sclient"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type IPGroupController struct {
	targetGroupARN string
}

var httpAddr = flag.String("http", ":8080", "Address to bind for metrics server")

func main() {
	flag.Parse()

	// Start metrics server
	http.Handle("/metrics", promhttp.Handler())
	go func() { log.Fatal(http.ListenAndServe(*httpAddr, nil)) }()

	// Setup k8s
	client, ns := client.NewK8sClient()

	core := client.Core()
	log.Println(core)
	pods, _ := core.Pods(ns).List(meta.ListOptions{})

	var unreadyPods unreadyPodList
	for _, p := range pods.Items {
		log.Println(p.Status.Conditions)
		unreadyPods = append(unreadyPods, &p)
	}
}

type unreadyPodList []*core.Pod
