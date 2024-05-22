package channelconsumer

import (
	"net"
	"time"
)

type Message struct {
	ID          int
	ChannelName string
	Content     interface{}
	ReceivedAt  time.Time
}

func NewMessage(id int, channelName string, content interface{}) *Message {
	return &Message{
		ID:          id,
		ChannelName: channelName,
		Content:     content,
		ReceivedAt:  time.Now(),
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
