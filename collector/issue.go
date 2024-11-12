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
	VideoIssuesMetric    *prometheus.Desc
	AudioIssuesMetric    *prometheus.Desc
	SubtitleIssuesMetric *prometheus.Desc
	OtherIssuesMetric    *prometheus.Desc
	OpenIssuesMetric     *prometheus.Desc
	ClosedIssuesMetric   *prometheus.Desc
}

func NewIssueCollector(client *goverseerr.Overseerr) *IssueCollector {
	logrus.Traceln("defining user collector")
	specificNamespace := "issue"
	return &IssueCollector{
		client: client,

		VideoIssuesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "video_total"),
			"Total number of video issues in overseerr",
			nil,
			prometheus.Labels{"url": client.URL, "type": "video"},
		),

		AudioIssuesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "audio_total"),
			"Total number of audio issues in overseerr",
			nil,
			prometheus.Labels{"url": client.URL},
		),

		SubtitleIssuesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "subtitle_total"),
			"Total number of subtitle issues in overseerr",
			nil,
			prometheus.Labels{"url": client.URL},
		),

		OtherIssuesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "other_total"),
			"Total number of other issues in overseerr",
			nil,
			prometheus.Labels{"url": client.URL},
		),

		OpenIssuesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "open_total"),
			"Total number of open issues in overseerr",
			nil,
			prometheus.Labels{"url": client.URL},
		),

		ClosedIssuesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "closed_total"),
			"Total number of closed issues in overseerr",
			nil,
			prometheus.Labels{"url": client.URL},
		),
	}
}

func (rc *IssueCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.VideoIssuesMetric
	ch <- rc.AudioIssuesMetric
	ch <- rc.SubtitleIssuesMetric
	ch <- rc.OtherIssuesMetric
	ch <- rc.OpenIssuesMetric
	ch <- rc.ClosedIssuesMetric
}

func (rc *IssueCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debugln("collecting issue data...")
	start := time.Now()
	result, err := rc.client.GetIssueCounts()
	if err != nil {
		logrus.WithField("error", err).Errorln("failed to get jobs from overseerr")
		return
	}

	ch <- prometheus.MustNewConstMetric(
		rc.VideoIssuesMetric,
		prometheus.GaugeValue,
		float64(result.Video),
	)

	ch <- prometheus.MustNewConstMetric(
		rc.AudioIssuesMetric,
		prometheus.GaugeValue,
		float64(result.Audio),
	)

	ch <- prometheus.MustNewConstMetric(
		rc.SubtitleIssuesMetric,
		prometheus.GaugeValue,
		float64(result.Subtitle),
	)

	ch <- prometheus.MustNewConstMetric(
		rc.OtherIssuesMetric,
		prometheus.GaugeValue,
		float64(result.Other),
	)

	ch <- prometheus.MustNewConstMetric(
		rc.OpenIssuesMetric,
		prometheus.GaugeValue,
		float64(result.Open),
	)

	ch <- prometheus.MustNewConstMetric(
		rc.ClosedIssuesMetric,
		prometheus.GaugeValue,
		float64(result.Closed),
	)

	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("job data collected")
}
