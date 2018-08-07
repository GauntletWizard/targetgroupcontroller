package main

import (
	"log"
	//"reflect"

	"github.com/prometheus/client_golang/prometheus"

	core "k8s.io/api/core/v1"
	watch "k8s.io/apimachinery/pkg/watch"
)

// Metrics
var (
	WatcherCommonLabels = []string{"Service"}
	eventsSeen          = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "Watcher",
			Name:      "events",
			Help:      "Number of Events seen by this Service",
		},
		WatcherCommonLabels,
	)
	channelClosed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "Watcher",
			Name:      "channelClosed",
			Help:      "Number of times the watch has closed",
		},
		WatcherCommonLabels,
	)
	readyIPs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "Watcher",
			Name:      "readyIPs",
			Help:      "Number of IPs that are ready from this service",
		},
		WatcherCommonLabels,
	)
	lastWatchEvent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "Watcher",
			Name:      "lastUpdated",
			Help:      "Last event from this service",
		},
		WatcherCommonLabels,
	)
)

func init() {
	prometheus.MustRegister(eventsSeen, readyIPs, lastWatchEvent)
}

// A Watcher observes changes to a kubernetes service (Endpoints) and updates a TargetGroup with the healthy hosts.
type Watcher struct {
	Service  string
	port     int64
	portName string
	tg       *TargetGroup
}

func (w *Watcher) GetReadyIPsFromEndpoint(e *core.Endpoints) (ips TargetSet) {
	ips = make(TargetSet)
	for _, subset := range e.Subsets {
		// Check if we're in the right port. Uncertain if a port can show up in multiple subsets, but we should deal with that possibility.
		var thisPort bool
		for _, port := range subset.Ports {
			if port.Name == w.portName || int64(port.Port) == w.port {
				thisPort = true
			}
		}
		thisPort = true
		if thisPort {
			for _, address := range subset.Addresses {
				ips[address.IP] = true
			}
		}
	}
	readyIPs.With(w.MetricLabels()).Set(float64(len(ips)))
	return
}

func (w *Watcher) Watch(events <-chan watch.Event) {
	for event := range events {
		endpoint := event.Object.(*core.Endpoints)
		if endpoint.Name == w.Service {
			eventsSeen.With(w.MetricLabels()).Inc()
			log.Println(endpoint.Name, w.GetReadyIPsFromEndpoint(endpoint))
			w.tg.Update(w.GetReadyIPsFromEndpoint(endpoint))
			lastWatchEvent.With(w.MetricLabels()).SetToCurrentTime()
		}
	}
	log.Println("Watch channel closed")
}

func (w *Watcher) MetricLabels() (l prometheus.Labels) {
	l = make(prometheus.Labels)
	l["Service"] = w.Service
	return
}
