package channelconsumer

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type MessageQueue interface {
	Add(message Message)
	Remove(message Message)
	SendPendingMessages(channelId int, connection net.Conn)
	Get() map[int][]Message
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
	fmt.Println("----------------------------------------------------------------------------------")
	logrus.Info("MessageCache after Added: ", mc.messages)
}

func (mc *InMemoryMessageCache) Remove(message Message) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	cacheMessages, found := mc.messages[message.ChannelId]

	if !found {
		return
	}

	for i, msg := range cacheMessages {
		if msg.ID == message.ID {
			mc.messages[message.ChannelId] = append(cacheMessages[:i], cacheMessages[i+1:]...)

			if len(mc.messages[message.ChannelId]) == 0 {
				delete(mc.messages, message.ChannelId)
			}

			logrus.Info("MessageCache after Deleted: ", mc.messages)
			break
		}
	}

}

func (mc *InMemoryMessageCache) SendPendingMessages(channelId int, connection net.Conn) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if messages, found := mc.messages[channelId]; found {
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

func (mc *InMemoryMessageCache) Get() map[int][]Message {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	return mc.messages
}
