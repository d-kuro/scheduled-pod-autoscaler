package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	minReplicasCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "scheduled_pod_auroscaler_min_replicas",
			Namespace: "scheduled_pod_auroscaler_controller",
			Help:      "Lower limit for the number of pods that can be set by the scheduled pod autoscaler",
		},
		[]string{"name", "namespace"},
	)

	maxReplicasCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "scheduled_pod_auroscaler_max_replicas",
			Namespace: "scheduled_pod_auroscaler_controller",
			Help:      "Upper limit for the number of pods that can be set by the scheduled pod autoscaler",
		},
		[]string{"name", "namespace"},
	)
)

func init() {
	metrics.Registry.MustRegister(minReplicasCounter, maxReplicasCounter)
}
