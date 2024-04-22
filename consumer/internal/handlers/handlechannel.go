package handlers

import (
	"fmt"

	"github.com/sagini18/message-broker/consumer/internal/types"
)

func HandleChannel(tcpConsumer *types.TcpConn) error {
	var input string
	fmt.Println("Enter the channel number you want to listen to: ")
	fmt.Scanln(&input)

	for input == "" {
		fmt.Println("Please provide a channel name. It can't be left blank.")
		fmt.Scanln(&input)
	}

	if _, err := tcpConsumer.Conn.Write([]byte(input)); err != nil {
		return err
	}
	return nil
}
