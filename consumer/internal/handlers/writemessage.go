package handlers

import (
	"fmt"
	"net"

	"github.com/sagini18/message-broker/consumer/internal/types"
	"github.com/sirupsen/logrus"
)

func WriteMessage(tcpConsumer net.Conn, receiver *types.Receiver) {
	if receiver.ReceivedMessage == nil {
		return
	}

	n, err := tcpConsumer.Write(receiver.ReceivedMessage)
	if err != nil {
		logrus.Error("Error in WriteMessage(): ", err)
		return
	}
	if n > 0 {
		fmt.Println("Sent message: ", string(receiver.ReceivedMessage))
	}

	receiver.ReceivedMessage = make([]byte, 0)
	receiver.ReadableReceivedMsgs = make([]types.Message, 0)
}
