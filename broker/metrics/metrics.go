package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ChannelGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "channel_count",
		Help: "Current number of channels",
	})
	RequestCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_count",
		Help: "Total number of requests served",
	})
)

func init() {
	prometheus.MustRegister(ChannelGauge, RequestCount)
}
