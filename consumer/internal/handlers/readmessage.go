package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sagini18/message-broker/consumer/internal/types"
	"github.com/sirupsen/logrus"
)

func ReadMessage(tcpConsumer *types.TcpConn, receiver *types.Receiver) {
	for {
		totalBytesRead := 0
		buffer := make([]byte, 200)

		for {
			n, err := tcpConsumer.Conn.Read(buffer[totalBytesRead:])
			if err != nil {
				logrus.Error("Error in reading data: ", err)
				return
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
		receiver.NewReceivedMessage(buffer[:totalBytesRead])

		unmarshalAndCallWriteMessage(receiver, tcpConsumer, totalBytesRead)
	}
}

func unmarshalAndCallWriteMessage(receiver *types.Receiver, tcpConsumer *types.TcpConn, totalBytesRead int) {
	if totalBytesRead > 0 {
		receivedData := make([]byte, totalBytesRead)
		copy(receivedData, receiver.ReceivedMessage[:totalBytesRead])

		chunks := bytes.Split(receivedData, []byte("]"))
		for _, chunk := range chunks {
			if len(chunk) > 0 {
				chunk = append(chunk, ']')

				error := json.Unmarshal(chunk, &receiver.ReadableReceivedMsgs)
				if error != nil {
					logrus.Error("Error in unmarshalling data: ", error)
					return
				}

				fmt.Println("------------------------------------------------------------------------------------------")

				if len(receiver.ReadableReceivedMsgs) > 0 {
					for _, msg := range receiver.ReadableReceivedMsgs {
						if len(receiver.ReadableReceivedMsgs) == 0 {
							continue
						}
						if msg.ChannelId == -1 && msg.MessageId == 0 {
							logrus.Info("Listening to channel: ", msg.Content)
							continue
						}

						logrus.Info("Received message: ", receiver.ReadableReceivedMsgs)

						receiver.ReceivedMessage = receivedData
						WriteMessage(tcpConsumer, receiver)
					}
				}
			}
		}
	}
}
