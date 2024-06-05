package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	// dto "github.com/prometheus/client_model/go"
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

// func GetChannelsEvents() int {
// 	var m dto.Metric
// 	ChannelsEvents.Write(&m)
// 	return int(m.GetGauge().GetValue())
// }

// func GetRequestsEvents() int {
// 	var m dto.Metric
// 	RequestsEvents.Write(&m)
// 	return int(m.GetCounter().GetValue())
// }
