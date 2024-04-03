package messagequeue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/internal/message"
)

var MessageCache message.CachedData = message.CachedData{Data: make(map[string][]message.Message)}

func AddToQueue(context echo.Context) error {
	id := context.Param("id")

	messageId := generateMessageId(id)

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

	if error := writeMessage(MessageCache.Data[id]); error != nil {
		fmt.Println("Error while writing message: ", error)
		return error
	}

	return context.JSON(http.StatusOK, MessageCache.Data[id])
}

func generateMessageId(id string) int {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	if len(MessageCache.Data[id]) == 0 {
		return 1
	}
	return MessageCache.Data[id][len(MessageCache.Data[id])-1].MessageId + 1
}

func writeMessage(messageCacheData []message.Message) error {
	messageBytes, err := json.Marshal(messageCacheData)
	if err != nil {
		return err
	}
	if err = message.Connection.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return err
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
}
