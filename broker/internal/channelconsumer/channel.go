package channelconsumer

import (
	"sync"
	"time"
)

type ChannelStorage interface {
	Add()
	Remove()
	Get() []ChannelEvent
}

type ChannelEvent struct {
	Timestamp time.Time
	Count     int
}

type Channel struct {
	mu           sync.Mutex
	channelEvent []ChannelEvent
	sseChannel   chan struct{}
}

func NewChannel() *Channel {
	return &Channel{
		channelEvent: []ChannelEvent{},
		sseChannel:   make(chan struct{}, 1),
	}
}

func (c *Channel) Add() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.channelEvent) < 1 {
		c.channelEvent = append(c.channelEvent, ChannelEvent{
			Timestamp: time.Now(),
			Count:     1,
		})
		c.notifySSE()
		return
	}

	lastCount := c.channelEvent[len(c.channelEvent)-1].Count
	c.channelEvent = append(c.channelEvent, ChannelEvent{
		Timestamp: time.Now(),
		Count:     lastCount + 1,
	},
	)
	c.notifySSE()
}

func (c *Channel) Remove() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.channelEvent) < 1 {
		return
	}

	lastCount := c.channelEvent[len(c.channelEvent)-1].Count
	c.channelEvent = append(c.channelEvent, ChannelEvent{
		Timestamp: time.Now(),
		Count:     lastCount - 1,
	},
	)
	c.notifySSE()
}

func (c *Channel) Get() []ChannelEvent {
	c.mu.Lock()
	defer c.mu.Unlock()

	return append([]ChannelEvent{}, c.channelEvent...)
}

func (c *Channel) SSEChannel() <-chan struct{} {
	return c.sseChannel
}

func (c *Channel) notifySSE() {
	select {
	case c.sseChannel <- struct{}{}:
	default:
	}
}
