package handlewebsocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/internal/message"
	"github.com/sagini18/message-broker/internal/messagequeue"
)

func HandleWebsocket(c echo.Context) error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	var error error

	message.Connection, error = upgrader.Upgrade(c.Response(), c.Request(), nil)
	if error != nil {
		return fmt.Errorf("UPGRADING_ERROR: %v", error)
	}
	defer message.Connection.Close()

	if error := readMessage(); error != nil {
		return fmt.Errorf("READING_ERROR: %v", error)
	}
	return nil
}

func readMessage() error {
	var msgs []message.Message
	for {
		_, receivedMessage, error := message.Connection.ReadMessage()
		if error != nil {
			return fmt.Errorf("READING_ERROR: %v", error)
		}

		error = json.Unmarshal(receivedMessage, &msgs)
		if error != nil {
			return fmt.Errorf("UNMARSHALING_ERROR: %v", error)
		}
		fmt.Println("Received message from consumer as ack: ", msgs)

		if err := messagequeue.RemoveMessageFromChannel(msgs); err != nil {
			return fmt.Errorf("REMOVING_MESSAGE_ERROR: %v", err)
		}
	}
}
