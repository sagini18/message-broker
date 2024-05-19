package channelconsumer

import (
	"net"
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
