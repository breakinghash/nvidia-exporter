package watch

import (
	"net/http"
	"strconv"
	"time"

	"github.com/uber-go/tally"

	log "github.com/sirupsen/logrus"
	promreporter "github.com/uber-go/tally/prometheus"
)

// Metrics represents temp and fan metrics
type Metrics struct {
	temp     map[uint]tally.Gauge
	fan      map[uint]tally.Gauge
	reporter promreporter.Reporter
}

// Init initializes temp and fan Gauge metrics
func (m Metrics) Init(prefix string, gpuCount uint) Metrics {
	m.reporter = promreporter.NewReporter(promreporter.Options{})

	rootScope, _ := tally.NewRootScope(tally.ScopeOptions{
		Prefix:         prefix,
		CachedReporter: m.reporter,
		Separator:      promreporter.DefaultSeparator,
	}, 30*time.Second)

	m.temp = make(map[uint]tally.Gauge)
	m.fan = make(map[uint]tally.Gauge)

	for i := uint(0); i < gpuCount; i++ {
		m.temp[i] = rootScope.Tagged(map[string]string{"gpu": strconv.Itoa(int(i))}).Gauge("temp")
		m.fan[i] = rootScope.Tagged(map[string]string{"gpu": strconv.Itoa(int(i))}).Gauge("fan")
	}

	return m
}

// ListenAndServe opens a prometheus metrics endpoint on 8080 port
func (m Metrics) ListenAndServe() {
	http.Handle("/metrics", m.reporter.HTTPHandler())

	log.Info("Serving :8080/metrics")
	log.Info("Stopped: %v", http.ListenAndServe(":8080", nil))
}
