package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/config"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/communication"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sagini18/message-broker/broker/internal/tcpconn"
	"github.com/sagini18/message-broker/broker/services/chart"
	"github.com/sagini18/message-broker/broker/services/table"
	"github.com/sirupsen/logrus"
)

func main() {
	configureLogger()

	go func() {
		logrus.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	config, err := config.LoadConfig()
	if err != nil {
		config.FilePath = "./internal/persistence/persisted_messages.txt"
	}
	file, err := os.OpenFile(config.FilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logrus.Error("Error in opening file: ", err)
	}
	defer file.Close()

	consumerStorage := channelconsumer.NewInMemoryInMemoryConsumerCache()
	messageQueue := channelconsumer.NewInMemoryMessageQueue()
	consumerIdGenerator := &channelconsumer.SerialConsumerIdGenerator{}
	messageIdGenerator := &channelconsumer.SerialMessageIdGenerator{}
	persist := persistence.New()
	requestCounter := channelconsumer.NewRequestCounter()
	failMsgCounter := channelconsumer.NewFailMsgCounter()
	tcpServer := tcpconn.New(":8081", consumerStorage, messageQueue, consumerIdGenerator, messageIdGenerator, persist, file)

	go func() {
		if err := tcpServer.Listen(); err != nil {
			logrus.Fatalf("Error in starting TCP server: %v", err)
		}
	}()

	app := echo.New()
	// app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{echo.GET, echo.POST},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	// }))

	// app.Use(middleware.Logger())

	app.POST("/api/channels/:id", func(c echo.Context) error {
		return communication.Broadcast(c, messageQueue, consumerStorage, messageIdGenerator, persist, file, requestCounter, failMsgCounter)
	})

	app.GET("/api/channel/all", func(c echo.Context) error {
		return table.ChannelDetails(c, messageQueue, consumerStorage, persist, file, requestCounter, failMsgCounter)
	})

	app.GET("/api/consumer/count", func(c echo.Context) error {
		return chart.Consumer(c, consumerStorage)
	})

	app.GET("/api/message/count", func(c echo.Context) error {
		return chart.Messages(c, messageQueue)
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
