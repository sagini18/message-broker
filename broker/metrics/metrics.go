package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ChannelsEvents = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "channels_events",
		Help: "Current number of channels",
	})
	RequestsEvents = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_events",
		Help: "Total number of requests served",
	})
	ConsumerEvents = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "consumers_events",
		Help: "Current number of consumers",
	})
	MessageEvents = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "messages_events",
		Help: "Current number of messages",
	})
)

func init() {
	prometheus.MustRegister(ChannelsEvents, RequestsEvents, ConsumerEvents, MessageEvents)
}
