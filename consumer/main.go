package main

import (
	"github.com/sagini18/consumers/internal/handlers"
)

func main() {
	handlers.InitConnection()

	handlers.ReadMessage()
}
