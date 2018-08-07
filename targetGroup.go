package main

import (
	"log"

	//	"github.com/aws/aws-sdk-go/aws"
	// 	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/prometheus/client_golang/prometheus"
)

var targetGroupARN string

var (
	allZones = "all"
)

// Metrics
var (
	TargetGroupCommonLabels = []string{"targetGroupARN"}
	targetsAdded            = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "TargetGroup",
			Name:      "targetsAdded",
			Help:      "Number of targets added to the TargetGroup",
		},
		TargetGroupCommonLabels,
	)
	targetsRemoved = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "TargetGroup",
			Name:      "targetsRemoved",
			Help:      "Number of targets removed from the TargetGroup",
		},
		TargetGroupCommonLabels,
	)
	healthyTargets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "TargetGroup",
			Name:      "healthyTargets",
			Help:      "Number of targets in TargetGroup marked as Healthy",
		},
		TargetGroupCommonLabels,
	)
	otherTargets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "TargetGroup",
			Name:      "otherTargets",
			Help:      "Number of targets in TargetGroup marked as something other than healthy",
		},
		TargetGroupCommonLabels,
	)
	registerErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "TargetGroup",
			Name:      "registerErrors",
			Help:      "Count of errors received from AWS when deregistering targets",
		},
		TargetGroupCommonLabels,
	)
	deregisterErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "TargetGroup",
			Name:      "deregisterErrors",
			Help:      "Count of errors received from AWS when deregistering targets",
		},
		TargetGroupCommonLabels,
	)
	describeErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "TargetGroup",
			Name:      "describeErrors",
			Help:      "Count of errors received from AWS when describing targetGroup",
		},
		TargetGroupCommonLabels,
	)
	targetGroupLastUpdated = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "TargetGroup",
			Name:      "lastUpdated",
			Help:      "Last finished updating targetGroup",
		},
		TargetGroupCommonLabels,
	)
)

func init() {
	prometheus.MustRegister(targetsAdded, targetsRemoved)
	prometheus.MustRegister(healthyTargets, otherTargets)
	prometheus.MustRegister(registerErrors, deregisterErrors, describeErrors)
	prometheus.MustRegister(targetGroupLastUpdated)
}

// A TargetGroup controls an ELBV2 TargetGroup and rectifies the list of endpoints as provided.
type TargetGroup struct {
	ARN               *string
	connection        *elbv2.ELBV2
	ExistingTargetSet TargetSet
	Port              int64
}

type TargetSet map[string]bool

func (tg *TargetGroup) Delta(ts TargetSet) (add []*elbv2.TargetDescription, remove []*elbv2.TargetDescription) {
	for t := range ts {
		if !tg.ExistingTargetSet[t] {
			// Make a copy of T because we need a fresh new pointer
			id := t
			log.Println("Adding:", id)
			add = append(add, &elbv2.TargetDescription{AvailabilityZone: &allZones,
				Id:   &id,
				Port: &tg.Port,
			})
		}
	}

	for t := range tg.ExistingTargetSet {
		if !ts[t] {
			id := t
			log.Println("Removing:", id)
			remove = append(remove, &elbv2.TargetDescription{AvailabilityZone: &allZones,
				Id:   &id,
				Port: &tg.Port,
			})
		}
	}
	return
}

func (tg *TargetGroup) Update(ts TargetSet) {
	new, old := tg.Delta(ts)

	log.Println("Beginning update of", tg.ARN)

	if len(new) != 0 {
		log.Println("Adding targets:", new)
		targetsAdded.With(tg.MetricLabels()).Add(float64(len(new)))
		_, err := tg.connection.RegisterTargets(&elbv2.RegisterTargetsInput{
			TargetGroupArn: tg.ARN,
			Targets:        new,
		})
		if err != nil {
			log.Println("Registration failed,", err)
			registerErrors.With(tg.MetricLabels()).Inc()
		}
	}

	if len(old) != 0 {
		log.Println("Removing targets:", old)
		targetsRemoved.With(tg.MetricLabels()).Add(float64(len(old)))
		_, err := tg.connection.DeregisterTargets(&elbv2.DeregisterTargetsInput{
			TargetGroupArn: tg.ARN,
			Targets:        old,
		})
		if err != nil {
			log.Println("Deregistration failed,", err)
			deregisterErrors.With(tg.MetricLabels()).Inc()
		}
	}
	tg.UpdateKnown()
	targetGroupLastUpdated.With(tg.MetricLabels()).SetToCurrentTime()
}

func (tg *TargetGroup) UpdateKnown() {
	targets, err := tg.connection.DescribeTargetHealth(
		&elbv2.DescribeTargetHealthInput{TargetGroupArn: tg.ARN})
	if err != nil {
		log.Println(targets, err)
		describeErrors.With(tg.MetricLabels()).Inc()
		return
	}

	current := make(TargetSet, len(targets.TargetHealthDescriptions))
	var healthy, unhealthy float64
	for _, t := range targets.TargetHealthDescriptions {
		current[*t.Target.Id] = true
		// Accumulate the healthcheck states of targets:
		if *t.TargetHealth.State == "healthy" {
			healthy += 1
		} else {
			unhealthy += 1
		}
	}
	healthyTargets.With(tg.MetricLabels()).Set(healthy)
	otherTargets.With(tg.MetricLabels()).Set(unhealthy)
	tg.ExistingTargetSet = current
}

func (tg TargetGroup) MetricLabels() (l prometheus.Labels) {
	l = make(prometheus.Labels)
	l["targetGroupARN"] = *tg.ARN
	return
}
