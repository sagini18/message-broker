package communication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
)

func Broadcast(context echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, messageIdGenerator *channelconsumer.SerialMessageIdGenerator, persist persistence.Persistence, file *os.File) error {
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

	if err := persist.Write(jsonBody, file); err != nil {
		logrus.Errorf("communication.Broadcast(): %v", err)
	}

	messageQueue.Add(*message)

	cachedMessages := messageQueue.Get(channelId)

	if err := writeMessage(cachedMessages, channelId, consumerStorage); err != nil { // run it in a background thread
		logrus.Errorf("communication.Broadcast(): %v", err)
		return context.JSON(http.StatusInternalServerError, "Failed to broadcast message")
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
			logrus.Errorf("writeMessage(): json.Marshal error: %v", err)
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
