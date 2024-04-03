package message

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	MessageId int
	ChannelId int
	Content   interface{}
}

type CachedData struct {
	sync.Mutex
	Data map[string][]Message
}

var Connection *websocket.Conn
