package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/camalot/goverseerr"
)


type SettingsCollector struct {
	client *goverseerr.Overseerr

	Settings *prometheus.Desc
}

func NewSettingsCollector(client *goverseerr.Overseerr) *SettingsCollector {
	logrus.Traceln("defining user collector")
	specificNamespace := "system"
	return &SettingsCollector{
		client: client,

		Settings: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "info"),
			"info of overseerr",
			[]string{"url", "appurl", "title"},
			nil,
		),
	}
}

func (rc *SettingsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Settings
}

func (rc *SettingsCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debugln("collecting settings data...")
	start := time.Now()
	settings, err := rc.client.GetMainSettings()
	status := 1
	if err != nil {
		logrus.WithField("error", err).Errorln("failed to get settings from overseerr")
		status = 0
	}
	if settings != nil {
		ch <- prometheus.MustNewConstMetric(
			rc.Settings,
			prometheus.GaugeValue,
			float64(status),
			rc.client.URL, settings.AppURL, settings.AppTitle,
		)
	}
	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("settings data collected")
}