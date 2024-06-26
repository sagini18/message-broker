package communication

import (
	"database/sql"
	"io"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/persistence"
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

func Broadcast(context echo.Context, messageQueue *channelconsumer.InMemoryMessageCache, consumerStorage *channelconsumer.InMemoryConsumerCache, messageIdGenerator *channelconsumer.SerialMessageIdGenerator, requestCount *channelconsumer.RequestCounter, failMsgCounter *channelconsumer.FailMsgCounter, channel *channelconsumer.Channel, database *sql.DB, persist persistence.Persistence) error {
	channelName := context.Param("id")
	requestCount.Add(channelName)

	if len(messageQueue.Get(channelName)) == 0 && len(consumerStorage.GetByChannel(channelName)) == 0 {
		channel.Add()
	}

	messageId := messageIdGenerator.NewId()
	message := channelconsumer.NewMessage(messageId, channelName, nil)
	context.Bind(message)

	go func() {
		if err := persist.Write(*message, database); err != nil {
			logrus.Errorf("communication.Broadcast(): persistence.WriteToDB() error: %v", err)
		}
	}()

	messageQueue.Add(*message)

	cachedMessages := messageQueue.Get(channelName)

	go func() {
		if error := writeMessage(cachedMessages, channelName, consumerStorage, messageQueue, channel); error != nil {
			failMsgCounter.Add(channelName)
			logrus.Errorf("communication.Broadcast(): writeMessage error: %v", error)
		}
	}()
	return context.JSON(http.StatusOK, message)
}

func writeMessage(messageCacheData []channelconsumer.Message, channelName string, store *channelconsumer.InMemoryConsumerCache, messageQueue *channelconsumer.InMemoryMessageCache, channel *channelconsumer.Channel) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
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
			if err == io.EOF || strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				if c := store.Get(consumer.Id, consumer.SubscribedChannel); c.TcpConn != nil {
					store.Remove(consumer.Id, consumer.SubscribedChannel)

					if len(store.GetByChannel(consumer.SubscribedChannel)) == 0 && len(messageQueue.Get(consumer.SubscribedChannel)) == 0 {
						channel.Remove()
					}
				}
				continue
			}
		}
	}
	return nil
}
