package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	// "github.com/willfantom/goverseerr"
	"github.com/camalot/goverseerr"
)


type JobCollector struct {
	client *goverseerr.Overseerr

	Job *prometheus.Desc
}

func NewJobCollector(client *goverseerr.Overseerr) *JobCollector {
	logrus.Traceln("defining user collector")
	specificNamespace := "job"
	return &JobCollector{
		client: client,

		Job: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "status"),
			"job status of overseerr",
			[]string{"url", "name", "id", "type"},
			nil,
		),
	}
}

func (rc *JobCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Job
}

func (rc *JobCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debugln("collecting job data...")
	start := time.Now()
	result, err := rc.client.GetJobs()
	status := 0
	if err != nil {
		logrus.WithField("error", err).Errorln("failed to get jobs from overseerr")
		status = 0
		return
	}
	for _, job := range result {
		if job == nil {
			continue
		}
		if job.Running {
			status = 1
		}
		ch <- prometheus.MustNewConstMetric(
			rc.Job,
			prometheus.GaugeValue,
			float64(status),
			rc.client.URL, job.Name, job.ID, string(job.Type),
		)
	}
	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("job data collected")
}