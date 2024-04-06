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

func (cd *CachedData) GenerateMessageId(id string) int {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	if len(cd.Data[id]) == 0 {
		return 1
	}
	return cd.Data[id][len(cd.Data[id])-1].MessageId + 1
}

type Consumer struct {
	ConsumerId         int
	SubscribedChannels []int
	TcpConn            net.Conn
}

type ConsumerCache struct {
	sync.Mutex
	Data []Consumer
}

func (ac *ConsumerCache) GenerateConsumerId() int {
	if len(ac.Data) == 0 {
		return 1
	}
	return ac.Data[len(ac.Data)-1].ConsumerId + 1
}

var ConsumerCacheData ConsumerCache = ConsumerCache{Data: []Consumer{}}
