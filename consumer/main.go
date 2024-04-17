package main

import (
	"github.com/sagini18/message-broker/consumer/internal/handlers"
)

func main() {
	handlers.InitConnection()

	handlers.ReadMessage()
}
