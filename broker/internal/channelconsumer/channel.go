package channelconsumer

import (
	"github.com/sagini18/message-broker/broker/metrics"
)

type ChannelStorage interface {
	Add()
	Remove()
}

type Channel struct {
	sseChannSum chan struct{}
}

func NewChannel() *Channel {
	return &Channel{
		sseChannSum: make(chan struct{}, 1),
	}
}

func (c *Channel) Add() {
	metrics.ChannelsEvents.Inc()
	c.notifySSE()
}

func (c *Channel) Remove() {
	metrics.ChannelsEvents.Dec()
	c.notifySSE()
}

func (c *Channel) SSEChannelSummary() <-chan struct{} {
	return c.sseChannSum
}

func (c *Channel) notifySSE() {
	select {
	case c.sseChannSum <- struct{}{}:
	default:
	}
}
