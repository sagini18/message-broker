package channel

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

/*
@Response
id: 9,
channelName: 76,
noOfMessagesInQueue: 12,
noOfConsumers: 65,
noOfProdcuers: 5,
nofMessagesInPersistence:3,
failedMessages: 1,
*/
func ChannelDetails(c echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache) error {
	messages := messageQueue.GetAll()
	consumers := consumerStorage.GetAll()
	if len(messages) == 0 && len(consumers) == 0 { //need to check with persisted file, producer count, failed messages count
		return c.JSON(http.StatusNoContent, "No data available")
	}

	var response []map[string]interface{}
	var channelInfo map[string]interface{}
	for channelId, messageList := range messages {
		channelInfo = make(map[string]interface{})
		channelInfo["channelName"] = channelId
		channelInfo["noOfMessagesInQueue"] = len(messageList)
		channelInfo["noOfConsumers"] = 0
		channelInfo["noOfProdcuers"] = 5
		channelInfo["nofMessagesInPersistence"] = 3
		channelInfo["failedMessages"] = 1
		response = append(response, channelInfo)
	}

	return c.JSON(http.StatusOK, response)
}
