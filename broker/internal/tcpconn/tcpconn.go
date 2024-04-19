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

	buf := make([]byte, 5120) //need to fix this
	for {
		n, err := connection.Read(buf)
		if err != nil {
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				store.Remove(consumer.Id)
				continue
			}
			return fmt.Errorf("tcpconn.listenToConsumerMessages(): connection.Read error: %v", err)

		}
		messageBytes := buf[:n]

		var msgs []channelconsumer.Message
		if err := json.Unmarshal(messageBytes, &msgs); err != nil {
			return fmt.Errorf("tcpconn.listenToConsumerMessages(): json.Unmarshal error: %v", err)

		}

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
