package channelconsumer

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type MessageStorage interface {
	Add(message Message)
	Remove(id int, channelName string)
	Get(channelName string) []Message
	SendPendingMessages(channelName string, connection net.Conn)
	GetAll() map[string][]Message
	GetCount(channelName string) int
}

type InMemoryMessageCache struct {
	mu       sync.RWMutex
	messages map[string][]Message
}

func NewInMemoryMessageQueue() *InMemoryMessageCache {
	return &InMemoryMessageCache{
		messages: make(map[string][]Message),
	}
}

func (mc *InMemoryMessageCache) Add(message Message) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if cacheMessages, found := mc.messages[message.ChannelName]; found {
		cacheMessages = append(cacheMessages, message)
		mc.messages[message.ChannelName] = cacheMessages
	} else {
		mc.messages[message.ChannelName] = []Message{message}
	}
	// logrus.Info("Added message to cache: ", message.Content)
}

func (mc *InMemoryMessageCache) Remove(id int, channelName string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	cacheMessages, found := mc.messages[channelName]
	if !found {
		return
	}

	var updatedMessages []Message
	for _, msg := range cacheMessages {
		if msg.ID != id {
			updatedMessages = append(updatedMessages, msg)
			continue
		}
		// logrus.Info("Removed message from cache: ", msg.Content)
	}
	mc.messages[channelName] = updatedMessages
	if len(updatedMessages) == 0 {
		delete(mc.messages, channelName)
	}
}

func (mc *InMemoryMessageCache) SendPendingMessages(channelName string, connection net.Conn) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	messagesCopy := make(map[string][]Message)
	for ChannelName, messages := range mc.messages {
		messagesCopy[ChannelName] = append([]Message{}, messages...)
	}

	if messages, found := messagesCopy[channelName]; found {
		messageBytes, err := json.Marshal(messages)
		if err != nil {
			logrus.Error("channelconsumer.SendPendingMessages() Error while marshalling message: ", err)
			return
		}

		if _, err = connection.Write(messageBytes); err != nil {
			logrus.Error("channelconsumer.SendPendingMessages() Error while writing previous messages to consumer: ", err)
			return
		}
	}
}

func (mc *InMemoryMessageCache) Get(channelName string) []Message {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	messagesCopy := append([]Message{}, mc.messages[channelName]...)
	return messagesCopy
}

func (mc *InMemoryMessageCache) GetAll() map[string][]Message {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	cacheMessagesCopy := make(map[string][]Message)
	for ChannelName, messages := range mc.messages {
		cacheMessagesCopy[ChannelName] = append([]Message{}, messages...)
	}
	return cacheMessagesCopy

}

func (mc *InMemoryMessageCache) GetCount(channelName string) int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return len(mc.messages[channelName])
}
