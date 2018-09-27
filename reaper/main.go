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
var dryrun = flag.Bool("dryrun", true, "Don't destroy pods, just print what you would destroy (default)")

const Ready = "Ready"

func main() {
	flag.Parse()

	// Start metrics server
	http.Handle("/metrics", promhttp.Handler())
	go func() { log.Fatal(http.ListenAndServe(*httpAddr, nil)) }()

	// Setup k8s
	client, ns := client.NewK8sClient()

	core := client.Core()
	pods, _ := core.Pods(ns).List(meta.ListOptions{})

	var unreadyPods unreadyPodList
	for _, p := range pods.Items {
		if podIsUnready(p) {
			unreadyPods = append(unreadyPods, &p)
		}
	}

}

// These consts describe  in our unready check
// (These should probably be annotated with official documentation)
const unreadyStatus = "False"
const completedReason = "PodCompleted"

func podIsUnready(p core.Pod) bool {
	for _, c := range p.Status.Conditions {
		if c.Type == Ready && c.Status == unreadyStatus && c.Reason != completedReason {
			return true
		}
	}
	return false
}

type unreadyPodList []*core.Pod
