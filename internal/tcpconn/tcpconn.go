package tcpconn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/sagini18/message-broker/internal/message"
	"github.com/sagini18/message-broker/internal/messagequeue"
)

func InitConnection() error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		return fmt.Errorf("LISTENING_ERROR: %v", err)
	}

	defer listener.Close()

	for {
		message.Connection, err = listener.Accept()
		if err != nil {
			return fmt.Errorf("ACCEPTING_ERROR: %v", err)
		}
		go handleClient(message.Connection)
	}
}

func handleClient(connection net.Conn) {
	consumer := message.Consumer{
		ConsumerId: message.ConsumerCacheData.GenerateConsumerId(),
		TcpConn:    connection,
	}

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

	message.ConsumerCacheData.Lock()
	message.ConsumerCacheData.Data = append(message.ConsumerCacheData.Data, consumer)
	message.ConsumerCacheData.Unlock()

	fmt.Println("Consumer subscribed to channel: ", message.ConsumerCacheData.Data)

	var channelConfir [1]message.Message
	channelConfir[0] = message.Message{
		MessageId: -1,
		ChannelId: -1,
		Content:   channel,
	}

	channelBytes, err := json.Marshal(channelConfir)
	if err != nil {
		fmt.Printf("MARSHALING_ERROR: %v\n", err)
		return
	}

	_, err = connection.Write(channelBytes)
	if err != nil {
		fmt.Printf("ERROR_WRITING_RESPONSE: %v\n", err)
		return
	}

	if err := readMessage(); err != nil {
		fmt.Println("Error while reading message: ", err)
	}
}

func readMessage() error {
	buf := make([]byte, 5120)
	for {
		for _, consumer := range message.ConsumerCacheData.Data {
			fmt.Println("-------------------------------------------------------------------")
			n, err := consumer.TcpConn.Read(buf)
			fmt.Println("Received message from consumer: ", string(buf[:n]))
			if err != nil {
				return fmt.Errorf("READING_ERROR: %v", err)
			}
			messageBytes := buf[:n]
			var msgs []message.Message

			chunks := bytes.Split(messageBytes, []byte("]"))
			for _, chunk := range chunks {
				if len(chunk) > 0 {
					chunk = append(chunk, ']')
					error := json.Unmarshal(chunk, &msgs)
					if error != nil {
						return fmt.Errorf("UNMARSHALING_ERROR: %v", error)
					}
					fmt.Println("Received message from consumer as ack: ", msgs)

					if err := messagequeue.RemoveMessageFromChannel(msgs); err != nil {
						return fmt.Errorf("REMOVING_MESSAGE_ERROR: %v", err)
					}
				}
			}
		}
	}
}
