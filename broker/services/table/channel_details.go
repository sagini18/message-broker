package table

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sagini18/message-broker/broker/sqlite"
	"github.com/sirupsen/logrus"
)

func ChannelDetails(c echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, persist persistence.Persistence, file *os.File, requestCounter *channelconsumer.RequestCounter, failMsgCount *channelconsumer.FailMsgCounter, channel *channelconsumer.Channel, sqlite sqlite.Persistence, database *sql.DB) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.String(http.StatusInternalServerError, "Streaming unsupported")
	}

	sendResponse := func() {
		response := channelSummaryResponse(messageQueue, consumerStorage, persist, file, requestCounter, failMsgCount, sqlite, database)
		data, err := json.Marshal(response)
		if err != nil {
			http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)
		flusher.Flush()
	}

	sendResponse()

	sseChannel := channel.SSEChannel()
	sseMessage := messageQueue.SSEChannel()
	sseConsumer := consumerStorage.SSEChannel()
	sseRequestCounter := requestCounter.SSEChannel()
	sseFailMsgCount := failMsgCount.SSEChannel()

	for {
		select {
		case <-sseChannel:
		case <-sseMessage:
		case <-sseConsumer:
		case <-sseRequestCounter:
		case <-sseFailMsgCount:
			sendResponse()
		case <-c.Request().Context().Done():
			return nil
		}
	}

}

func channelSummaryResponse(messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, persist persistence.Persistence, file *os.File, requestCounter *channelconsumer.RequestCounter, failMsgCount *channelconsumer.FailMsgCounter, sqlite sqlite.Persistence, database *sql.DB) []map[string]interface{} {
	messages := messageQueue.GetAll()
	consumers := consumerStorage.GetAll()
	persistMessages, err := persist.ReadAll(file)
	if err != nil {
		logrus.Error("ChannelDetails(): error reading from persistence file: ", err)
	}
	dbmsgs, err := sqlite.ReadAll(database)
	if err != nil {
		logrus.Error("ChannelDetails(): error reading from persistence db: ", err)
	}
	fmt.Println("dbmsgs: ", dbmsgs)
	failedChannels := failMsgCount.GetAllChannel()

	if len(messages) == 0 && len(consumers) == 0 && len(persistMessages) == 0 && len(failedChannels) == 0 {
		return []map[string]interface{}{}
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
	id := 1
	for channelName := range channelSet {
		count := sqlite.ReadCount(channelName, database)
		fmt.Println("Count: ", count)
		channelInfo := map[string]interface{}{
			"id":                        id,
			"channelName":               channelName,
			"noOfMessagesInQueue":       messageQueue.GetCount(channelName),
			"noOfConsumers":             len(consumers[channelName]),
			"noOfRequests":              requestCounter.Get(channelName),
			"noOfMessagesInPersistence": persist.ReadCount(channelName, file),
			"failedMessages":            failMsgCount.Get(channelName),
		}
		response = append(response, channelInfo)
		id++
	}
	return response
}
