package handlers

import (
	"fmt"
	"net"

	"github.com/sagini18/consumers/internal/types"
)

func InitConnection() {
	var err error
	types.Connection, err = net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println(err)
		return
	}

	var input string
	fmt.Println("Enter the channel number you want to listen to: ")
	fmt.Scanln(&input)

	for input == "" {
		fmt.Println("Please provide a channel name. It can't be left blank.")
		fmt.Scanln(&input)
	}

	if _, err = types.Connection.Write([]byte(input)); err != nil {
		fmt.Println("Error sending channel number:", err)
		return
	}
}
