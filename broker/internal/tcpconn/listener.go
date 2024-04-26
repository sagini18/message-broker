package tcpconn

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

var messageChan = make(chan channelconsumer.Message)

func listenToConsumerMessages(connection net.Conn, consumer *channelconsumer.Consumer, store channelconsumer.Storage) error {
	defer connection.Close()

	for {
		buffer, totalBytesRead, err := readMessages(connection, store, consumer)
		if err != nil {
			return fmt.Errorf("tcpconn.listenToConsumerMessages(): connection.Read error: %v", err)
		}

		messageBytes := buffer[:totalBytesRead]

		var msgs []channelconsumer.Message
		if err := json.Unmarshal(messageBytes, &msgs); err != nil {
			return fmt.Errorf("tcpconn.listenToConsumerMessages(): json.Unmarshal error: %v", err)
		}

		fmt.Println(" msgs", msgs)

		for _, msg := range msgs {
			messageChan <- msg
		}
	}
}

func removeMessages(queue channelconsumer.MessageQueue) {
	for msg := range messageChan {
		logrus.Info("Received message from consumer as ack: ", msg)

		queue.Remove(msg)
	}
}

func readMessages(connection net.Conn, store channelconsumer.Storage, consumer *channelconsumer.Consumer) ([]byte, int, error) {
	totalBytesRead := 0
	buffer := make([]byte, 200)

	for {
		n, err := connection.Read(buffer[totalBytesRead:])
		if err != nil {
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				if c := store.GetConsumer(consumer.Id); c.Id == consumer.Id {
					store.Remove(consumer.Id)
				}
				continue
			}
			return buffer, totalBytesRead, err
		}

		totalBytesRead += n

		if totalBytesRead >= len(buffer) {
			newBufferSize := len(buffer) * 2
			newBuffer := make([]byte, newBufferSize)
			copy(newBuffer, buffer)
			buffer = newBuffer
		} else {
			break
		}
	}
	return buffer, totalBytesRead, nil
}
