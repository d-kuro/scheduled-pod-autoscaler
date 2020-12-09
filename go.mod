module github.com/d-kuro/scheduled-pod-autoscaler

go 1.15

require (
	github.com/go-logr/logr v0.1.0
	github.com/google/go-cmp v0.5.3
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.4
	github.com/prometheus/client_golang v1.8.0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.4
)
