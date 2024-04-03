package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/internal/handlewebsocket"
	"github.com/sagini18/message-broker/internal/messagequeue"
)

func main() {
	app := echo.New()

	app.GET("/ws", handlewebsocket.HandleWebsocket)

	app.POST("/api/channels/:id", messagequeue.AddToQueue)

	app.Logger.Fatal(app.Start(":8080"))
}
