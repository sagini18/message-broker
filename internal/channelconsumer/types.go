package channelconsumer

import (
	"net"
	"sync"
)

type Message struct {
	MessageId int
	ChannelId int
	Content   interface{}
}

type CachedData struct {
	sync.Mutex
	Data map[int][]Message
}

func (cd *CachedData) generateMessageId(id int) int {
	cd.Lock()
	defer cd.Unlock()

	if id == -1{
		return -1
	}else if len(cd.Data[id]) == 0 {
		return 1
	}
	return cd.Data[id][len(cd.Data[id])-1].MessageId + 1
}

func NewMessage(channelNum int, content interface{}) *Message {	
	return &Message{
		MessageId: MessageCache.generateMessageId(channelNum),
		ChannelId: channelNum,
		Content:   content,
	}
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

func (ac *ConsumerCache) generateConsumerId() int {
	ac.Lock()
	defer ac.Unlock()

	if len(ac.Data) == 0 {
		return 1
	}
	return ac.Data[len(ac.Data)-1].ConsumerId + 1
}

func NewConsumer(conn *net.Conn) *Consumer {
	return &Consumer{
		ConsumerId: ConsumerCacheData.generateConsumerId(),
		TcpConn:    *conn,
	}
}

var ConsumerCacheData ConsumerCache = ConsumerCache{Data: []Consumer{}}
var MessageCache CachedData = CachedData{Data: make(map[int][]Message)}
