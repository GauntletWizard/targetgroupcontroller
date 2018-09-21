package main

import (
	"flag"
	"log"
	"net/http"

	// 	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	client, _ := client.NewK8sClient()

	for {
		core := client.Core()
		log.Println(core)
	}
}
