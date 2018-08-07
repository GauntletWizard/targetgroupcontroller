package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type IPGroupController struct {
	targetGroupARN string
}

var region = flag.String("region", "", "AWS Region")
var targetARN = flag.String("targetgrouparn", "", "ARN of the Target group to update")
var service = flag.String("service", "", "Name of the service to watch")
var port = flag.Int64("port", 0, "Port number")
var httpAddr = flag.String("http", ":8080", "Address to bind for metrics server")

func main() {
	flag.Parse()

	if *region == "" || *targetARN == "" || *service == "" || *port == 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	// Start metrics server
	http.Handle("/metrics", promhttp.Handler())
	go func() { log.Fatal(http.ListenAndServe(*httpAddr, nil)) }()

	// Setup AWS
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	elbSvc := elbv2.New(sess, &aws.Config{Region: region})

	tg := TargetGroup{
		ARN:        targetARN,
		connection: elbSvc,
		Port:       9000,
	}
	tg.UpdateKnown()

	watcher := Watcher{Service: *service, port: *port, tg: &tg}

	// Setup k8s
	client, ns := NewK8sClient()
	watchOptions := meta.ListOptions{}

	for {
		watch, err := client.Core().Endpoints(ns).Watch(watchOptions)
		if err != nil {
			log.Fatal("Opening watch failed,", err)
		}
		watchChan := watch.ResultChan()
		watcher.Watch(watchChan)
	}
}
