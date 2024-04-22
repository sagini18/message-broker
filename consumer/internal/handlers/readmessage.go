package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sagini18/message-broker/consumer/internal/types"
)

func ReadMessage() {
	for {
		totalBytesRead := 0
		buffer := make([]byte, 200)

		for {
			n, err := types.Connection.Read(buffer[totalBytesRead:])
			if err != nil {
				fmt.Println("Error in reading data: ", err)
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

		types.ReceivedMessage = buffer[:totalBytesRead]

		fmt.Println("------------------------------------------------------------------------------------------")

		if totalBytesRead > 0 {
			receivedData := make([]byte, totalBytesRead)
			copy(receivedData, types.ReceivedMessage[:totalBytesRead])

			chunks := bytes.Split(receivedData, []byte("]"))
			for _, chunk := range chunks {
				if len(chunk) > 0 {
					chunk = append(chunk, ']')

					error := json.Unmarshal(chunk, &types.ReadableReceivedMsgs)
					if error != nil {
						fmt.Println("Error in unmarshalling: ", error)
						return
					}

					if len(types.ReadableReceivedMsgs) > 0 {
						for _, msg := range types.ReadableReceivedMsgs {
							if len(types.ReadableReceivedMsgs) == 0 {
								continue
							}
							if msg.ChannelId == -1 && msg.MessageId == 0 {
								fmt.Println("Listening to the channel: ", msg.Content)
								continue
							}

							fmt.Println("Received : ", types.ReadableReceivedMsgs)

							types.ReceivedMessage = receivedData
							WriteMessage()
						}
					}
				}
			}
		}
	}
}
