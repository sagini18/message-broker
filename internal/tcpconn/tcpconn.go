package tcpconn

import (
	"encoding/json"
	"fmt"
	"net"
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

	message.Connection, err = listener.Accept()
	if err != nil {
		return fmt.Errorf("ACCEPTING_ERROR: %v", err)
	}

	readMessage(message.Connection)
	return nil
}

func readMessage(connection net.Conn) error {
	defer connection.Close()

	buf := make([]byte, 1024)
	for {
		n, err := connection.Read(buf)
		if err != nil {
			return fmt.Errorf("READING_ERROR: %v", err)
		}
		messageBytes := buf[:n]
		var msgs []message.Message

		error := json.Unmarshal(messageBytes, &msgs)
		if error != nil {
			return fmt.Errorf("UNMARSHALING_ERROR: %v", error)
		}
		fmt.Println("Received message from consumer as ack: ", msgs)

		if err := messagequeue.RemoveMessageFromChannel(msgs); err != nil {
			return fmt.Errorf("REMOVING_MESSAGE_ERROR: %v", err)
		}
	}
}
