package messagequeue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/internal/message"
	"golang.org/x/exp/slices"
)

func AddToQueue(context echo.Context) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	id := context.Param("id")

	messageId := message.MessageCache.GenerateMessageId(id)

	channelId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error while converting id to int: ", err)
		return err
	}
	messageBody := message.Message{
		MessageId: messageId,
		ChannelId: channelId,
	}
	context.Bind(&messageBody)

	saveMessageToCache(id, messageBody)

	if error := writeMessage(message.MessageCache.Data[id], channelId); error != nil {
		fmt.Println("Error while writing message: ", error)
		return error
	}

	return context.JSON(http.StatusOK, message.MessageCache.Data[id])
}

func writeMessage(messageCacheData []message.Message, id int) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	for _, consumer := range message.ConsumerCacheData.Data {
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

func saveMessageToCache(id string, messageBody message.Message) {
	message.MessageCache.Lock()
	defer message.MessageCache.Unlock()

	if cachedData, found := message.MessageCache.Data[id]; found {
		cachedData = append(cachedData, messageBody)
		message.MessageCache.Data[id] = cachedData
	} else {
		message.MessageCache.Data[id] = []message.Message{messageBody}

	}

	fmt.Println("Message saved to cache: ", message.MessageCache.Data)
	if id == "1" {
		fmt.Println("-------------------------------------------------------------------")
	}
}
