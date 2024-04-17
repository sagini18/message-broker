package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sagini18/message-broker/consumer/internal/types"
)

func ReadMessage() {
	for {
		n, err := types.Connection.Read(types.ReceivedMessage)
		if err != nil {
			fmt.Println("Error in reading: ", err)
			return
		}
		fmt.Println("-------------------------------------------------------------")

		if n > 0 {
			receivedData := make([]byte, n)
			copy(receivedData, types.ReceivedMessage[:n])

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
							if msg.ChannelId == -1 && msg.MessageId == -1 {
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
