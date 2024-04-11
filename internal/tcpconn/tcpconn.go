package tcpconn

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/sagini18/message-broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/internal/messagequeue"
)

func InitConnection() error {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		return fmt.Errorf("LISTENING_ERROR: %v", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("ACCEPTING_ERROR: %v", err)
		}
		go handleClient(conn)
	}
}

var messageChan = make(chan channelconsumer.Message)

func handleClient(connection net.Conn) {
	consumer := channelconsumer.NewConsumer(&connection)

	channelBuf := make([]byte, 5120)
	n, err := connection.Read(channelBuf)
	if err != nil {
		fmt.Printf("ERROR_READING_CHANNEL: %v", err)
		return
	}
	channel := string(channelBuf[:n])
	channelInt, err := strconv.Atoi(channel)
	if err != nil {
		fmt.Printf("CONVERSION_ERROR: %v\n", err)
		return
	}

	consumer.SubscribedChannels = append(consumer.SubscribedChannels, channelInt)

	channelconsumer.ConsumerCacheData.Lock()
	channelconsumer.ConsumerCacheData.Data = append(channelconsumer.ConsumerCacheData.Data, *consumer)
	channelconsumer.ConsumerCacheData.Unlock()

	fmt.Println("Consumer subscribed to channel: ", channelconsumer.ConsumerCacheData.Data)

	var channelConfirmation [1]channelconsumer.Message
	channelConfirmation[0] = *channelconsumer.NewMessage(-1, channel)

	channelBytes, err := json.Marshal(channelConfirmation)
	if err != nil {
		fmt.Printf("MARSHALING_ERROR: %v\n", err)
		return
	}

	_, err = connection.Write(channelBytes)
	if err != nil {
		fmt.Printf("ERROR_WRITING_RESPONSE: %v\n", err)
		return
	}

	sendPreviousMessages(channelInt, connection)

	go func() {
		defer connection.Close()

		buf := make([]byte, 5120)
		for {
			n, err := connection.Read(buf)
			if err != nil {
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
	}()

	removeMessages()
}

func removeMessages() {
	for msg := range messageChan {
		fmt.Println("Received message from consumer as ack: ", msg)

		messagequeue.RemoveMessageFromChannel(msg)
	}
}

func sendPreviousMessages(channelId int, connection net.Conn) {
	channelconsumer.MessageCache.Lock()
	defer channelconsumer.MessageCache.Unlock()

	if messages, found := channelconsumer.MessageCache.Data[channelId]; found {
		messageBytes, err := json.Marshal(messages)
		if err != nil {
			fmt.Println("Error while marshalling message: ", err)
			return
		}

		_, err = connection.Write(messageBytes)
		if err != nil {
			fmt.Println("Error while writing previous messages to consumer: ", err)
			return
		}
	}
}
