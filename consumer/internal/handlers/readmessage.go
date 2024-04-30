package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	"github.com/sagini18/message-broker/consumer/internal/types"
	"github.com/sirupsen/logrus"
)

func ReadMessage(tcpConsumer net.Conn, receiver *types.Receiver) {
	for {

		buffer, totalBytesRead, err := readAndExpandBuffer(tcpConsumer)
		if err != nil {
			logrus.Error("Error in reading data: ", err)
			return
		}

		receiver.NewReceivedMessage(buffer[:totalBytesRead])

		unmarshalMessage(receiver, tcpConsumer, totalBytesRead)
	}
}

func readAndExpandBuffer(tcpConsumer net.Conn) ([]byte, int, error) {
	totalBytesRead := 0
	buffer := make([]byte, 200)

	for {
		n, err := tcpConsumer.Read(buffer[totalBytesRead:])
		if err != nil {
			return buffer, totalBytesRead, err
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

func unmarshalMessage(receiver *types.Receiver, tcpConsumer net.Conn, totalBytesRead int) {
	if totalBytesRead <= 0 {
		return
	}
	receiver.ReceivedMessage = receiver.ReceivedMessage[:totalBytesRead]

	chunks := bytes.Split(receiver.ReceivedMessage, []byte("]"))
	for _, chunk := range chunks {
		if len(chunk) <= 0 {
			continue
		}
		chunk = append(chunk, ']')

		if err := json.Unmarshal(chunk, &receiver.ReadableReceivedMsgs); err != nil {
			logrus.Error("Error in unmarshalling data: ", err)
			return
		}

		fmt.Println("------------------------------------------------------------------------------------------")

		decodeMessage(tcpConsumer, receiver)
	}
}

func decodeMessage(tcpConsumer net.Conn, receiver *types.Receiver) {
	if len(receiver.ReadableReceivedMsgs) <= 0 {
		return
	}
	for _, msg := range receiver.ReadableReceivedMsgs {
		if len(receiver.ReadableReceivedMsgs) == 0 {
			continue
		}
		if msg.ChannelId == -1 {
			logrus.Info("Listening to channel: ", msg.Content)
			continue
		}

		logrus.Info("Received message: ", receiver.ReadableReceivedMsgs)

		WriteMessage(tcpConsumer, receiver)
	}
}
