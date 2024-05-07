package channelconsumer

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type MessageStorage interface {
	Add(message Message)
	Remove(message Message)
	SendPendingMessages(channelId int, connection net.Conn)
	GetMessages(channelId int) []Message
}

type InMemoryMessageCache struct {
	mu       sync.Mutex
	messages map[int][]Message
}

func NewInMemoryMessageStore() *InMemoryMessageCache {
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
	logrus.Info("MessageCache after Added: ", mc.messages)
}

func (mc *InMemoryMessageCache) Remove(message Message) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	cacheMessages, found := mc.messages[message.ChannelId]
	if !found {
		return
	}

	var updatedMessages []Message
	for _, msg := range cacheMessages {
		if msg.ID != message.ID {
			updatedMessages = append(updatedMessages, msg)
		}
	}
	mc.messages[message.ChannelId] = updatedMessages
	if len(updatedMessages) == 0 {
		delete(mc.messages, message.ChannelId)
	}
	logrus.Info("MessageCache after Removed: ", mc.messages)
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

func (mc *InMemoryMessageCache) GetMessages(channelId int) []Message {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	return mc.messages[channelId]
}
