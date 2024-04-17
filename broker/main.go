package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/messagequeue"
	"github.com/sagini18/message-broker/broker/internal/tcpconn"
)

func main() {

	app := echo.New()

	go func() {
		if err := tcpconn.InitConnection(); err != nil {
			app.Logger.Fatal(err)
		}
	}()

	app.POST("/api/channels/:id", messagequeue.AddToQueue)

	app.Logger.Fatal(app.Start(":8080"))
}
