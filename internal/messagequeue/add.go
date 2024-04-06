package messagequeue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/internal/message"
)

var MessageCache message.CachedData = message.CachedData{Data: make(map[string][]message.Message)}

func AddToQueue(context echo.Context) error {
	id := context.Param("id")

	messageId := MessageCache.GenerateMessageId(id)

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

	if messageBody.ChannelId == 10 {
		if error := writeMessage(MessageCache.Data[id]); error != nil {
			fmt.Println("Error while writing message: ", error)
			return error
		}
	}

	return context.JSON(http.StatusOK, MessageCache.Data[id])
}

func writeMessage(messageCacheData []message.Message) error {
	for _, consumer := range message.ConsumerCacheData.Data {
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
	return nil
}

func saveMessageToCache(id string, messageBody message.Message) {
	MessageCache.Lock()
	defer MessageCache.Unlock()

	if cachedData, found := MessageCache.Data[id]; found {
		cachedData = append(cachedData, messageBody)
		MessageCache.Data[id] = cachedData
	}
	MessageCache.Data[id] = []message.Message{messageBody}

	fmt.Println("Message saved to cache: ", MessageCache.Data)
	if id == "1" {
		fmt.Println("-------------------------------------------------------------------")
	}
}
