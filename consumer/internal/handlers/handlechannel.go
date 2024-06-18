package handlers

import (
	"fmt"
	"net"
	"os"
)

func HandleChannel(tcpConsumer net.Conn) error {
	var input string
	fmt.Println("Enter the channel number you want to listen to: ")
	fmt.Scanln(&input)

	if input == "" {
		input = os.Getenv("CHANNEL_NAME")
		for input == "" {
			fmt.Println("Please provide a channel name. It can't be left blank.")
			fmt.Scanln(&input)
		}
	}

	if _, err := tcpConsumer.Write([]byte(input)); err != nil {
		return err
	}
	return nil
}
