package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/goverseerr"
)


type CacheCollector struct {
	client *goverseerr.Overseerr

	Cache *prometheus.Desc
}

func NewCacheCollector(client *goverseerr.Overseerr) *CacheCollector {
	logrus.Traceln("defining cache collector")
	return &CacheCollector{
		client: client,

		Cache: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "cache"),
			"cache status of overseerr",
			[]string{"url", "name", "id", "type"},
			nil,
		),
	}
}

func (rc *CacheCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Cache
}

func (rc *CacheCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debugln("collecting cache data...")
	start := time.Now()
	result, err := rc.client.GetCacheStats()
	if err != nil {
		logrus.WithField("error", err).Errorln("failed to get cache from overseerr")
		return
	}
	for _, cache := range result {
		if cache == nil {
			continue
		}
		
		// add hits
		ch <- prometheus.MustNewConstMetric(
			rc.Cache,
			prometheus.GaugeValue,
			float64(cache.Stats.Hits),
			rc.client.URL, cache.Name, cache.ID, "hits",
		)

		// add misses
		ch <- prometheus.MustNewConstMetric(
			rc.Cache,
			prometheus.GaugeValue,
			float64(cache.Stats.Misses),
			rc.client.URL, cache.Name, cache.ID, "misses",
		)

		// add key count
		ch <- prometheus.MustNewConstMetric(
			rc.Cache,
			prometheus.GaugeValue,
			float64(cache.Stats.Keys),
			rc.client.URL, cache.Name, cache.ID, "keys",
		)

		// add value size
		ch <- prometheus.MustNewConstMetric(
			rc.Cache,
			prometheus.GaugeValue,
			float64(cache.Stats.VSize),
			rc.client.URL, cache.Name, cache.ID, "vsize",
		)

		// add key size
		ch <- prometheus.MustNewConstMetric(
			rc.Cache,
			prometheus.GaugeValue,
			float64(cache.Stats.KSize),
			rc.client.URL, cache.Name, cache.ID, "ksize",
		)
	}
	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("cache data collected")
}