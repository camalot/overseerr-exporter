package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/goverseerr"
)


type StatusCollector struct {
	client *goverseerr.Overseerr

	Status *prometheus.Desc
}

func NewStatusCollector(client *goverseerr.Overseerr) *StatusCollector {
	logrus.Traceln("defining user collector")
	specificNamespace := "system"
	return &StatusCollector{
		client: client,

		Status: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "status"),
			"status of overseerr",
			[]string{"url", "version", "commit_tag"},
			nil,
		),
	}
}

func (rc *StatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Status
}

func (rc *StatusCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debugln("collecting status data...")
	start := time.Now()
	result, err := rc.client.Status()
	status := 1
	if err != nil {
		logrus.WithField("error", err).Errorln("failed to get status from overseerr")
		status = 0
		return
	}
	if result != nil {
		
		ch <- prometheus.MustNewConstMetric(
			rc.Status,
			prometheus.GaugeValue,
			float64(status),
			rc.client.URL, result.Version, result.CommitTag,
		)
	}
	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("status data collected")
}