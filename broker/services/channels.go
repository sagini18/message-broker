package services

import (
	"database/sql"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/persistence"
	"github.com/sirupsen/logrus"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Channels(c echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, requestCounter *channelconsumer.RequestCounter, failMsgCount *channelconsumer.FailMsgCounter, channel *channelconsumer.Channel, persist persistence.Persistence, database *sql.DB) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.JSON(http.StatusNotImplemented, map[string]string{
			"type":    "StreamError",
			"message": "Streaming unsupported",
		})
	}

	sendResponse := func() {
		response := channelSummaryResponse(messageQueue, consumerStorage, requestCounter, failMsgCount, persist, database)
		data, err := json.Marshal(response)
		if err != nil {
			http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)
		flusher.Flush()
	}

	sendResponse()

	sseChannel := channel.SSEChannelSummary()
	sseMessage := messageQueue.SSEChannelSummary()
	sseConsumer := consumerStorage.SSEChannelSummary()
	sseRequestCounter := requestCounter.SSEChannelSummary()
	sseFailMsgCount := failMsgCount.SSEChannelSummary()

	for {
		select {
		case <-sseChannel:
			sendResponse()
		case <-sseMessage:
			sendResponse()
		case <-sseConsumer:
			sendResponse()
		case <-sseRequestCounter:
			sendResponse()
		case <-sseFailMsgCount:
			sendResponse()
		case <-c.Request().Context().Done():
			return nil
		}
	}

}

func channelSummaryResponse(messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, requestCounter *channelconsumer.RequestCounter, failMsgCount *channelconsumer.FailMsgCounter, persist persistence.Persistence, database *sql.DB) []map[string]interface{} {
	messages := messageQueue.GetAll()
	consumers := consumerStorage.GetAll()
	dbmsgs, err := persist.ReadAll(database)
	if err != nil {
		logrus.Error("ChannelDetails(): error reading from persistence db: ", err)
	}
	failedChannels := failMsgCount.GetAllChannel()

	if len(messages) == 0 && len(consumers) == 0 && len(dbmsgs) == 0 && len(failedChannels) == 0 {
		return []map[string]interface{}{}
	}

	if len(dbmsgs) > len(messages) {
		messages = dbmsgs
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
	id := 1
	for channelName := range channelSet {
		channelInfo := map[string]interface{}{
			"id":                        id,
			"channelName":               channelName,
			"noOfMessagesInQueue":       messageQueue.GetCount(channelName),
			"noOfConsumers":             len(consumers[channelName]),
			"noOfRequests":              requestCounter.Get(channelName),
			"noOfMessagesInPersistence": persist.ReadCount(channelName, database),
			"failedMessages":            failMsgCount.Get(channelName),
		}
		response = append(response, channelInfo)
		id++
	}
	return response
}
