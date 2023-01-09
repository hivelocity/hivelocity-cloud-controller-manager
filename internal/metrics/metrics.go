/*
TODO: boilerplate
*/

package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

var (
	OperationCalled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cloud_controller_manager_operations_total",
		Help: "The total number of operation was called",
	}, []string{"op"})
)

var registry = prometheus.NewRegistry()

func GetRegistry() *prometheus.Registry {
	return registry
}

func Serve(address string) {
	klog.Info("Starting metrics server at ", address)

	registry.MustRegister(OperationCalled)

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}

	http.Handle("/metrics", promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(address, nil); err != nil {
		klog.ErrorS(err, "create metrics service")
	}
}
