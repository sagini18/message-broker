package tcpconn

import (
	"bytes"
	"database/sql"
	"fmt"
	"net"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/persistence"
	"github.com/sirupsen/logrus"
)

func listenToConsumerMessages(connection net.Conn, consumer *channelconsumer.Consumer, store channelconsumer.Storage, messageQueue channelconsumer.MessageStorage, channel channelconsumer.ChannelStorage, perist persistence.Persistence, database *sql.DB) error {
	defer connection.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	for {
		buffer, totalBytesRead, err := readMessages(connection, store, consumer, messageQueue, channel)
		if err != nil {
			return fmt.Errorf("listenToConsumerMessages(): connection.Read error: %v", err)
		}

		if totalBytesRead <= 0 {
			continue
		}
		messageBytes := buffer[:totalBytesRead]

		var msgs []channelconsumer.Message
		chunks := bytes.Split(messageBytes, []byte("]"))

		for _, chunk := range chunks {
			trimedChunk := bytes.TrimFunc(chunk, func(r rune) bool {
				return r == 0
			})

			if len(trimedChunk) <= 0 {
				continue
			}
			trimedChunk = append(trimedChunk, ']')

			if err := json.Unmarshal(trimedChunk, &msgs); err != nil {
				logrus.Errorf("listenToConsumerMessages(): json.Unmarshal error: %v", err)
				continue
			}
		}

		for _, msg := range msgs {
			logrus.Info("Message received as ack: ", msg)

			go func() {
				if err := perist.Remove(msg.ID, database); err != nil {
					logrus.Errorf("tcpconn.listenToConsumerMessages(): persistence.RemoveFromDB() error: %v", err)
				}
			}()
			messageQueue.Remove(msg.ID, msg.ChannelName)

			if messageQueue.GetCount(msg.ChannelName) == 0 && len(store.GetByChannel(msg.ChannelName)) == 0 {
				channel.Remove()
			}
		}
	}
}

func readMessages(connection net.Conn, store channelconsumer.Storage, consumer *channelconsumer.Consumer, messageQueue channelconsumer.MessageStorage, channel channelconsumer.ChannelStorage) ([]byte, int, error) {
	totalBytesRead := 0
	buffer := make([]byte, 200)

	for {
		n, err := connection.Read(buffer[totalBytesRead:])
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				if c := store.Get(consumer.Id, consumer.SubscribedChannel); c.TcpConn != nil {
					store.Remove(consumer.Id, consumer.SubscribedChannel)

					if len(store.GetByChannel(consumer.SubscribedChannel)) == 0 && len(messageQueue.Get(consumer.SubscribedChannel)) == 0 {
						channel.Remove()
					}
				}
				continue
			}
		}

		totalBytesRead += n

		if totalBytesRead >= len(buffer) {
			newBufferSize := len(buffer) * 2
			newBuffer := make([]byte, newBufferSize)
			copy(newBuffer, buffer)
			buffer = newBuffer
			continue
		}

		return buffer, totalBytesRead, nil
	}
}
