package tcpconn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
)

func listenToConsumerMessages(connection net.Conn, consumer *channelconsumer.Consumer, store channelconsumer.Storage, messageQueue channelconsumer.MessageStorage, persist persistence.Persistence, file *os.File) error {
	defer connection.Close()

	for {
		buffer, totalBytesRead, err := readMessages(connection, store, consumer)
		if err != nil {
			return fmt.Errorf("tcpconn.listenToConsumerMessages(): connection.Read error: %v", err)
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
				logrus.Errorf("tcpconn.listenToConsumerMessages(): json.Unmarshal error: %v", err)
				continue
			}
		}

		for _, msg := range msgs {
			logrus.Info("Message received as ack: ", msg)

			if err := persist.Remove(msg.ID, file); err != nil {
				logrus.Errorf("tcpconn.listenToConsumerMessages(): persistence.Remove() error: %v", err)
			}

			messageQueue.Remove(msg.ID, msg.ChannelId)
		}
	}
}

func readMessages(connection net.Conn, store channelconsumer.Storage, consumer *channelconsumer.Consumer) ([]byte, int, error) {
	totalBytesRead := 0
	buffer := make([]byte, 200)

	for {
		n, err := connection.Read(buffer[totalBytesRead:])
		if err != nil {
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				if c := store.Get(consumer.Id); c.TcpConn != nil {
					store.Remove(consumer.Id)
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
