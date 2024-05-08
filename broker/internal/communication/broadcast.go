package communication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
)

func Broadcast(context echo.Context, messageStore *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, messageIdGenerator *channelconsumer.SerialMessageIdGenerator) error {
	id := context.Param("id")

	channelId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("communication.Broadcast(): strconv.Atoi error: %v", err)
	}

	lastMessageId, err := persistence.GetLastMessageId(channelId)
	if err != nil {
		logrus.Errorf("communication.Broadcast(): getLastMessageId error: %v", err)
	}
	messageIdGenerator.SetLastId(lastMessageId)

	newId := messageIdGenerator.NewId()
	messageBody := channelconsumer.NewMessage(newId, channelId, nil)
	context.Bind(messageBody)

	jsonBody, err := json.Marshal(messageBody)
	if err != nil {
		return fmt.Errorf("communication.Broadcast(): json.Marshal error: %v", err)
	}

	if err := persistence.AppendToFile(jsonBody); err != nil {
		logrus.Errorf("communication.Broadcast(): persistence.AppendToFile error: %v", err)
	}

	messageStore.Add(*messageBody)

	messageCacheData := messageStore.GetMessages(channelId)

	if error := writeMessage(messageCacheData, channelId, consumerStorage); error != nil {
		logrus.Errorf("communication.Broadcast(): writeMessage error: %v", error)
	}

	return context.JSON(http.StatusOK, messageCacheData)
}

func writeMessage(messageCacheData []channelconsumer.Message, id int, store *channelconsumer.InMemoryConsumerCache) error {
	allConsumers := store.Get()

	for _, consumer := range allConsumers {
		if !slices.Contains(consumer.SubscribedChannels, id) {
			continue
		}

		messageBytes, err := json.Marshal(messageCacheData)
		if err != nil {
			logrus.Errorf("communication.writeMessage(): json.Marshal error: %v", err)
			continue
		}

		if _, err := consumer.TcpConn.Write(messageBytes); err != nil {
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				if c := store.GetConsumer(consumer.Id); c.TcpConn != nil {
					store.Remove(consumer.Id)
				}
				continue
			}

		}

	}
	return nil
}
