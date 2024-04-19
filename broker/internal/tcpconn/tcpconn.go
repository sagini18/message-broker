package tcpconn

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

var messageChan = make(chan channelconsumer.Message)

func listenToConsumerMessages(connection net.Conn, consumer *channelconsumer.Consumer, store channelconsumer.Storage) {

	buf := make([]byte, 5120)
	for {
		n, err := connection.Read(buf)
		if err != nil {
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				store.Remove(consumer.Id)
				continue
			}
			fmt.Printf("READING_ERROR: %v", err)
			return
		}
		messageBytes := buf[:n]

		var msgs []channelconsumer.Message
		if err := json.Unmarshal(messageBytes, &msgs); err != nil {
			fmt.Printf("UNMARSHALING_ERROR: %v", err)
			return
		}

		for _, msg := range msgs {
			messageChan <- msg
		}
	}

}

func removeMessages(queue channelconsumer.MessageQueue) {
	for msg := range messageChan {
		fmt.Println("Received message from consumer as ack: ", msg)

		queue.Remove(msg)
	}
}
