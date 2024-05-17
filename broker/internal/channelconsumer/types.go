package channelconsumer

import (
	"net"
)

type Message struct {
	ID        int
	ChannelId int //need to change into string
	Content   interface{}
}

func NewMessage(id int, channelNum int, content interface{}) *Message {
	return &Message{
		ID:        id,
		ChannelId: channelNum,
		Content:   content,
	}
}

type Consumer struct {
	Id                 int
	SubscribedChannels []int
	TcpConn            net.Conn
}

func NewConsumer(id int, conn net.Conn, subscribedChannels []int) *Consumer {
	return &Consumer{
		Id:                 id,
		TcpConn:            conn,
		SubscribedChannels: subscribedChannels,
	}
}
