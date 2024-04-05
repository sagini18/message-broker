package message

import (
	"net"
	"sync"
)

var Connection net.Conn

type Message struct {
	MessageId int
	ChannelId int
	Content   interface{}
}

type CachedData struct {
	sync.Mutex
	Data map[string][]Message
}
