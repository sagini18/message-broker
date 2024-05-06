package communication

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

func Broadcast(context echo.Context, messageStore *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, messageIdGenerator *channelconsumer.SerialMessageIdGenerator) error {
	id := context.Param("id")

	channelId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("communication.Broadcast(): strconv.Atoi error: %v", err)
	}

	newId := messageIdGenerator.NewId()
	messageBody := channelconsumer.NewMessage(newId, channelId, nil)
	context.Bind(messageBody)

	messageStore.Add(*messageBody)

	allMessages := messageStore.Get()
	messageCacheData := allMessages[channelId]

	if error := writeMessage(messageCacheData, channelId, consumerStorage); error != nil {
		logrus.Errorf("communication.Broadcast(): writeMessage error: %v", error)
	}

	return context.JSON(http.StatusOK, allMessages[channelId])
}

func writeMessage(messageCacheData []channelconsumer.Message, id int, store *channelconsumer.InMemoryConsumerCache) error {
	allConsumers := store.Get()

	for _, consumer := range allConsumers {
		if !slices.Contains(consumer.SubscribedChannels, id) {
			continue
		}

		messageBytes, err := json.Marshal(messageCacheData)
		if err != nil {
			return fmt.Errorf("communication.writeMessage(): json.Marshal error: %v", err)
		}

		if _, err := consumer.TcpConn.Write(messageBytes); err != nil {
			if opErr, ok := err.(*net.OpError); !ok && opErr.Op != "write" {
				continue
			}

			store.Remove(consumer.Id)
			continue
		}

	}
	return nil
}
