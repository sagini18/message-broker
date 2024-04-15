package types

import (
	"net"
)

type Message struct {
	MessageId int
	ChannelId int
	Content   interface{}
}

var Connection net.Conn
var ReceivedMessage = make([]byte, 5120)
var ReadableReceivedMsgs []Message
