package messagequeue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func AddToQueue(context echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, messageIdGenerator *channelconsumer.SerialMessageIdGenerator) error {
	id := context.Param("id")

	channelId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error while converting id to int: ", err)
		return err
	}

	newId := messageIdGenerator.NewId()
	messageBody := channelconsumer.NewMessage(newId, channelId, nil)
	context.Bind(messageBody)

	messageQueue.Add(*messageBody)

	allMessages := messageQueue.Get()

	if error := writeMessage(allMessages[channelId], channelId, consumerStorage); error != nil {
		fmt.Println("Error while writing message: ", error)
		return error
	}

	return context.JSON(http.StatusOK, allMessages[channelId])
}

func writeMessage(messageCacheData []channelconsumer.Message, id int, store *channelconsumer.InMemoryConsumerCache) error {
	allConsumers := store.Get()

	for _, consumer := range allConsumers {
		if slices.Contains(consumer.SubscribedChannels, id) {
			messageBytes, err := json.Marshal(messageCacheData)
			if err != nil {
				fmt.Println("Error while marshalling message: ", err)
				return err
			}

			if _, err := consumer.TcpConn.Write(messageBytes); err != nil {
				if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
					store.Remove(consumer.Id)
					continue
				}
				fmt.Printf("WRITING_ERROR: %v", err.Error())
				return err
			}
		}
	}
	return nil
}
