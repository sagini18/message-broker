package table

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
)

func ChannelDetails(c echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, persist persistence.Persistence, file *os.File, requestCounter *channelconsumer.RequestCounter, failMsgCount *channelconsumer.FailMsgCounter) error {
	messages := messageQueue.GetAll()
	consumers := consumerStorage.GetAll()
	persistMessages, err := persist.ReadAll(file)
	if err != nil {
		logrus.Error("ChannelDetails(): error reading from persistence file: ", err)
	}
	failedChannels := failMsgCount.GetAllChannel()

	if len(messages) == 0 && len(consumers) == 0 && len(persistMessages) == 0 && len(failedChannels) == 0 {
		return c.JSON(http.StatusNoContent, "No data available")
	}

	if len(persistMessages) > len(messages) {
		messages = persistMessages
	}

	channelSet := make(map[string]struct{})
	for channelName := range messages {
		channelSet[channelName] = struct{}{}
	}
	for channelName := range consumers {
		channelSet[channelName] = struct{}{}
	}
	for _, channelName := range failedChannels {
		channelSet[channelName] = struct{}{}
	}

	response := make([]map[string]interface{}, 0, len(channelSet))
	for channelName := range channelSet {
		channelInfo := map[string]interface{}{
			"channelName":               channelName,
			"noOfMessagesInQueue":       messageQueue.GetCount(channelName),
			"noOfConsumers":             len(consumers[channelName]),
			"noOfRequests":              requestCounter.Get(channelName),
			"noOfMessagesInPersistence": persist.ReadCount(channelName, file),
			"failedMessages":            failMsgCount.Get(channelName),
		}
		response = append(response, channelInfo)
	}
	return c.JSON(http.StatusOK, response)
}
