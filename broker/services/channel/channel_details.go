package channel

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
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
func ChannelDetails(c echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, persist persistence.Persistence, file *os.File, producerCount *channelconsumer.ProducerCounter, failMsgCount *channelconsumer.FailMsgCounter) error {
	messages := messageQueue.GetAll()
	consumers := consumerStorage.GetAll()
	persistMessages, err := persist.ReadAll(file)
	if err != nil {
		logrus.Error("ChannelDetails(): error reading from persistence file: ", err)
	}

	if len(messages) == 0 && len(consumers) == 0 && len(persistMessages) == 0 { //need to check producer count, failed messages count
		return c.JSON(http.StatusNoContent, "No data available")
	}

	if len(persistMessages) > len(messages) {
		messages = persistMessages
	}

	var response []map[string]interface{}
	var channelInfo map[string]interface{}
	for channelName := range messages { //need to check consumers  //it is only printing when there is no consumers for this channel
		channelInfo = make(map[string]interface{})
		channelInfo["channelName"] = channelName
		channelInfo["noOfMessagesInQueue"] = messageQueue.GetCount(channelName)
		channelInfo["noOfConsumers"] = len(consumers[channelName])
		channelInfo["noOfProdcuers"] = producerCount.Get(channelName)
		channelInfo["nofMessagesInPersistence"] = persist.ReadCount(channelName, file)
		channelInfo["failedMessages"] = failMsgCount.Get(channelName)
		response = append(response, channelInfo)
	}
	return c.JSON(http.StatusOK, response)
}
