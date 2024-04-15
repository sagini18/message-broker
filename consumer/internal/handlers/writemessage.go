package handlers

import (
	"fmt"
	"log"

	"github.com/sagini18/consumers/internal/types"
)

func WriteMessage() {

	if types.ReceivedMessage != nil {
		_, err := types.Connection.Write(types.ReceivedMessage)

		fmt.Println("WriteMessage: ", string(types.ReceivedMessage))

		if err != nil {
			log.Println("write:", err)
			return
		}
		types.ReceivedMessage = make([]byte, 1024)
		types.ReadableReceivedMsgs = []types.Message{}
	}

}
