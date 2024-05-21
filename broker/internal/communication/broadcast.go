package communication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
)

/*
* Broadcast handles the broadcasting of messages to consumers on a specific channel.
* It is responsible for processing the incoming HTTP request, creating a new message,
* persisting it, adding it to the message queue, and then sending it to all relevant consumers.
*
* Parameters:
*   - context (echo.Context): The Echo context containing the HTTP request and response.
*   - messageQueue (*channelconsumer.InMemoryMessageQueue): The in-memory queue where messages are stored before being sent to consumers.
*   - consumerStorage (*channelconsumer.InMemoryConsumerCache): The in-memory cache of active consumers that are listening for messages.
*   - messageIdGenerator (*channelconsumer.SerialMessageIdGenerator): The generator for creating unique message IDs.
*   - persist (persistence.Persistence): The persistence layer responsible for saving messages to a file.
*   - file (*os.File): The file where messages are persisted.
*
* Returns:
*   - error: Returns an error if there is any issue during the process, otherwise returns nil.
 */

func Broadcast(context echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, messageIdGenerator *channelconsumer.SerialMessageIdGenerator, persist persistence.Persistence, file *os.File, producerCount *channelconsumer.ProducerCounter, failMsgCounter *channelconsumer.FailMsgCounter) error {
	channelName := context.Param("id")
	producerCount.Add(channelName)

	messageId := messageIdGenerator.NewId()
	message := channelconsumer.NewMessage(messageId, channelName, nil)
	context.Bind(message)

	jsonBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("communication.Broadcast(): json.Marshal error: %v", err)
	}
	go func() {
		if err := persist.Write(jsonBody, file); err != nil {
			logrus.Errorf("communication.Broadcast(): persistence.Write() error: %v", err)
		}
	}()

	messageQueue.Add(*message)

	cachedMessages := messageQueue.Get(channelName)

	go func() {
		if error := writeMessage(cachedMessages, channelName, consumerStorage); error != nil {
			failMsgCounter.Add(channelName)
			logrus.Errorf("communication.Broadcast(): writeMessage error: %v", error)
		}
	}()
	return context.JSON(http.StatusOK, cachedMessages)
}

func writeMessage(messageCacheData []channelconsumer.Message, channelName string, store *channelconsumer.InMemoryConsumerCache) error {
	consumers := store.GetByChannel(channelName)

	if len(consumers) == 0 {
		return nil
	}

	for _, consumer := range consumers {
		messageBytes, err := json.Marshal(messageCacheData)
		if err != nil {
			logrus.Errorf("writeMessage(): json.Marshal error: %v", err)
			continue
		}

		if _, err := consumer.TcpConn.Write(messageBytes); err != nil {
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				if c := store.Get(consumer.Id, consumer.SubscribedChannel); c.TcpConn != nil {
					store.Remove(consumer.Id, consumer.SubscribedChannel)
				}
				continue
			}
		}
	}
	return nil
}
