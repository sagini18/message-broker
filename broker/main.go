package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/communication"
	"github.com/sagini18/message-broker/broker/internal/tcpconn"
	"github.com/sirupsen/logrus"
)

func main() {
	configureLogger()

	consumerIdGenerator := &channelconsumer.SerialConsumerIdGenerator{}
	consumerStorage := channelconsumer.NewInMemoryInMemoryConsumerCache()
	messageIdGenerator := &channelconsumer.SerialMessageIdGenerator{}
	messageQueue := channelconsumer.NewInMemoryMessageQueue()
	tcpServer := tcpconn.New(":8081", consumerStorage, messageQueue, consumerIdGenerator, messageIdGenerator)

	go func() {
		if err := tcpServer.Listen(); err != nil {
			logrus.Fatalf("Error in starting TCP server: %v", err)
		}
	}()

	app := echo.New()
	app.POST("/api/channels/:id", func(c echo.Context) error {
		return communication.Broadcast(c, messageQueue, consumerStorage, messageIdGenerator)
	})
	if err := app.Start(":8080"); err != nil {
		logrus.Fatalf("Error in starting API server: %v", err)
	}
}

func configureLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
