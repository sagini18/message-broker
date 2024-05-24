package channelconsumer

import (
	"net"
	"sync"
	"time"
)

type Message struct {
	ID          int
	ChannelName string
	Content     interface{}
}

func NewMessage(id int, channelName string, content interface{}) *Message {
	return &Message{
		ID:          id,
		ChannelName: channelName,
		Content:     content,
	}
}

type Consumer struct {
	Id                int
	SubscribedChannel string
	TcpConn           net.Conn
}

func NewConsumer(id int, conn net.Conn, subscribedChannel string) *Consumer {
	return &Consumer{
		Id:                id,
		TcpConn:           conn,
		SubscribedChannel: subscribedChannel,
	}
}

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
}

func NewChannel() *Channel {
	return &Channel{
		channelEvent: []ChannelEvent{},
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
		return
	}

	lastCount := c.channelEvent[len(c.channelEvent)-1].Count
	c.channelEvent = append(c.channelEvent, ChannelEvent{
		Timestamp: time.Now(),
		Count:     lastCount + 1,
	},
	)
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
}

func (c *Channel) Get() []ChannelEvent {
	c.mu.Lock()
	defer c.mu.Unlock()

	return append([]ChannelEvent{}, c.channelEvent...)
}
