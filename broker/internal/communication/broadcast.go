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
	"github.com/sirupsen/logrus"
)

func Broadcast(context echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, messageIdGenerator *channelconsumer.SerialMessageIdGenerator) error {
	id := context.Param("id")

	channelId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("messagequeue.AddToQueue(): strconv.Atoi error: %v", err)
	}

	newId := messageIdGenerator.NewId()
	messageBody := channelconsumer.NewMessage(newId, channelId, nil)
	context.Bind(messageBody)

	messageQueue.Add(*messageBody)

	allMessages := messageQueue.Get()

	if error := writeMessage(allMessages[channelId], channelId, consumerStorage); error != nil {
		logrus.Errorf("messagequeue.AddToQueue(): writeMessage error: %v", error)
		return context.JSON(http.StatusInternalServerError, "Error in writing message to consumer: "+error.Error())
	}

	return context.JSON(http.StatusOK, allMessages[channelId])
}

func writeMessage(messageCacheData []channelconsumer.Message, id int, store *channelconsumer.InMemoryConsumerCache) error {
	allConsumers := store.Get()

	for _, consumer := range allConsumers {
		if slices.Contains(consumer.SubscribedChannels, id) {
			messageBytes, err := json.Marshal(messageCacheData)
			if err != nil {
				return fmt.Errorf("messagequeue.writeMessage(): json.Marshal error: %v", err)
			}

			if _, err := consumer.TcpConn.Write(messageBytes); err != nil {
				if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
					store.Remove(consumer.Id)
					continue
				}
				return fmt.Errorf("messagequeue.writeMessage(): consumer.TcpConn.Write error: %v", err)
			}
		}
	}
	return nil
}
