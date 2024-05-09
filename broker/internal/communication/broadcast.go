package communication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
)

func Broadcast(context echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, messageIdGenerator *channelconsumer.SerialMessageIdGenerator, persist persistence.Persistence) error {

	channelId, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return fmt.Errorf("communication.Broadcast(): strconv.Atoi error: %v", err)
	}

	messageId := messageIdGenerator.NewId()
	message := channelconsumer.NewMessage(messageId, channelId, nil)
	context.Bind(message)

	jsonBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("communication.Broadcast(): json.Marshal error: %v", err)
	}

	time.Sleep(10 * time.Second)
	if err := persist.Write(jsonBody); err != nil {
		logrus.Errorf("communication.Broadcast(): persistence.AppendToFile error: %v", err)
	}

	messageQueue.Add(*message)

	cachedMessages := messageQueue.Get(channelId)

	if error := writeMessage(cachedMessages, channelId, consumerStorage); error != nil {
		logrus.Errorf("communication.Broadcast(): writeMessage error: %v", error)
	}

	return context.JSON(http.StatusOK, cachedMessages)
}

func writeMessage(messageCacheData []channelconsumer.Message, id int, store *channelconsumer.InMemoryConsumerCache) error {
	allConsumers := store.GetAll()

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
				if c := store.Get(consumer.Id); c.TcpConn != nil {
					store.Remove(consumer.Id)
				}
				continue
			}

		}

	}
	return nil
}
