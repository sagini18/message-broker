package channelconsumer

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type MessageStorage interface {
	Add(message Message)
	Remove(id int, channelId int)
	Get(channelId int) []Message
	SendPendingMessages(channelId int, connection net.Conn)
	GetAll() map[int][]Message
}

type InMemoryMessageCache struct {
	mu       sync.Mutex
	messages map[int][]Message
}

func NewInMemoryMessageQueue() *InMemoryMessageCache {
	return &InMemoryMessageCache{
		messages: make(map[int][]Message),
	}
}

func (mc *InMemoryMessageCache) Add(message Message) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if cacheMessages, found := mc.messages[message.ChannelId]; found {
		cacheMessages = append(cacheMessages, message)
		mc.messages[message.ChannelId] = cacheMessages
	} else {
		mc.messages[message.ChannelId] = []Message{message}
	}
	logrus.Info("Added message to cache: ", message.Content)
}

func (mc *InMemoryMessageCache) Remove(id int, channelId int) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	cacheMessages, found := mc.messages[channelId]
	if !found {
		return
	}

	var updatedMessages []Message
	for _, msg := range cacheMessages {
		if msg.ID != id {
			updatedMessages = append(updatedMessages, msg)
			continue
		}
		logrus.Info("Removed message from cache: ", msg.Content)
	}
	mc.messages[channelId] = updatedMessages
	if len(updatedMessages) == 0 {
		delete(mc.messages, channelId)
	}
}

func (mc *InMemoryMessageCache) SendPendingMessages(channelId int, connection net.Conn) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	messagesCopy := make(map[int][]Message)
	for k, v := range mc.messages {
		messagesCopy[k] = append([]Message{}, v...)
	}

	if messages, found := messagesCopy[channelId]; found {
		messageBytes, err := json.Marshal(messages)
		if err != nil {
			logrus.Error("SendPendingMessages() Error while marshalling message: ", err)
			return
		}

		if _, err = connection.Write(messageBytes); err != nil {
			logrus.Error("SendPendingMessages() Error while writing previous messages to consumer: ", err)
			return
		}
	}
}

func (mc *InMemoryMessageCache) Get(channelId int) []Message {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	return mc.messages[channelId]
}

func (mc *InMemoryMessageCache) GetAll() map[int][]Message {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	return mc.messages
}
