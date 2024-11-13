package collector

import (
	"time"

	"github.com/camalot/goverseerr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type IssueCollector struct {
	client *goverseerr.Overseerr

	// IssueMetric       *prometheus.Desc
	IssueTypeMetric   *prometheus.Desc
	IssueStatusMetric *prometheus.Desc
}

func NewIssueCollector(client *goverseerr.Overseerr) *IssueCollector {
	logrus.Traceln("defining user collector")
	specificNamespace := "issue"
	return &IssueCollector{
		client: client,

		IssueTypeMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "type"),
			"Total number of issues in overseerr by type",
			[]string{"type"},
			prometheus.Labels{"url": client.URL},
		),

		IssueStatusMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "status"),
			"Total number of issues in overseerr by status",
			[]string{"status"},
			prometheus.Labels{"url": client.URL},
		),
	}
}

func (rc *IssueCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.IssueTypeMetric
	ch <- rc.IssueStatusMetric
}

func (rc *IssueCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debugln("collecting issue data...")
	start := time.Now()
	result, err := rc.client.GetIssueCounts()
	if err != nil {
		logrus.WithField("error", err).Errorln("failed to get issue counts from overseerr")
		return
	}

	var types = map[string]int{
		"video":     result.Video,
		"audio":     result.Audio,
		"subtitle":  result.Subtitle,
		"other":     result.Other,
	}

	var statuses = map[string]int{
		"open":   result.Open,
		"closed": result.Closed,
	}

	for k, v := range types {
		ch <- prometheus.MustNewConstMetric(
			rc.IssueTypeMetric,
			prometheus.GaugeValue,
			float64(v),
			k,
		)
	}

	for k, v := range statuses {
		ch <- prometheus.MustNewConstMetric(
			rc.IssueStatusMetric,
			prometheus.GaugeValue,
			float64(v),
			k,
		)
	}

	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("issue count data collected")
}
