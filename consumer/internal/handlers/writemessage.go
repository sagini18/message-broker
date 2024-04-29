package handlers

import (
	"fmt"

	"github.com/sagini18/message-broker/consumer/internal/types"
	"github.com/sirupsen/logrus"
)

func WriteMessage(tcpConsumer *types.TcpConn, receiver *types.Receiver) {
	if receiver.ReceivedMessage == nil {
		return
	}

	if _, err := tcpConsumer.Conn.Write(receiver.ReceivedMessage); err != nil {
		logrus.Error("Error in writing data: ", err)
		return
	}
	fmt.Println("Sent message: ", string(receiver.ReceivedMessage))

	receiver.ReceivedMessage = make([]byte, 1024)
	receiver.ReadableReceivedMsgs = make([]types.Message, 0)
}
