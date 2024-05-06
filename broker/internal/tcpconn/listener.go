package tcpconn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

func listenToConsumerMessages(connection net.Conn, consumer *channelconsumer.Consumer, store channelconsumer.Storage, messageStore channelconsumer.MessageStorage) error {
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
				return fmt.Errorf("tcpconn.listenToConsumerMessages(): json.Unmarshal error: %v", err)
			}
		}

		for _, msg := range msgs {
			logrus.Info("Message received as ack: ", msg)
			messageStore.Remove(msg)
		}
	}
}

func readMessages(connection net.Conn, store channelconsumer.Storage, consumer *channelconsumer.Consumer) ([]byte, int, error) {
	totalBytesRead := 0
	buffer := make([]byte, 200)

	for {
		n, err := connection.Read(buffer[totalBytesRead:])
		if err != nil {
			if opErr, ok := err.(*net.OpError); !ok && opErr.Op != "read" { //race conditions only in the image
				continue
			}
			if c := store.GetConsumer(consumer.Id); c.TcpConn != nil {
				store.Remove(consumer.Id)
			}
			continue
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
