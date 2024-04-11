package messagequeue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/internal/channelconsumer"
	"golang.org/x/exp/slices"
)

func AddToQueue(context echo.Context) error {
	id := context.Param("id")

	channelId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error while converting id to int: ", err)
		return err
	}
	messageBody := channelconsumer.NewMessage(channelId, nil)
	context.Bind(messageBody)

	saveMessageToCache(channelId, *messageBody)

	if error := writeMessage(channelconsumer.MessageCache.Data[channelId], channelId); error != nil {
		fmt.Println("Error while writing message: ", error)
		return error
	}

	return context.JSON(http.StatusOK, channelconsumer.MessageCache.Data[channelId])
}

func writeMessage(messageCacheData []channelconsumer.Message, id int) error {
	for _, consumer := range channelconsumer.ConsumerCacheData.Data {
		if slices.Contains(consumer.SubscribedChannels, id) {
			messageBytes, err := json.Marshal(messageCacheData)
			if err != nil {
				fmt.Println("Error while marshalling message: ", err)
				return err
			}

			if _, err := consumer.TcpConn.Write(messageBytes); err != nil {
				fmt.Println("Error while writing message to consumer: ", err)
				return err
			}
		}
	}
	return nil
}

func saveMessageToCache(id int, messageBody channelconsumer.Message) {
	channelconsumer.MessageCache.Lock()
	defer channelconsumer.MessageCache.Unlock()

	if cachedData, found := channelconsumer.MessageCache.Data[id]; found {
		cachedData = append(cachedData, messageBody)
		channelconsumer.MessageCache.Data[id] = cachedData
	} else {
		channelconsumer.MessageCache.Data[id] = []channelconsumer.Message{messageBody}

	}

	fmt.Println("Message saved to cache: ", channelconsumer.MessageCache.Data)
}
