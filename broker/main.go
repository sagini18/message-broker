package main

import (
	"database/sql"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sagini18/message-broker/broker/config"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/communication"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sagini18/message-broker/broker/internal/tcpconn"
	"github.com/sagini18/message-broker/broker/services/chart"
	"github.com/sagini18/message-broker/broker/services/table"
	"github.com/sagini18/message-broker/broker/sqlite"
	"github.com/sirupsen/logrus"
)

func main() {
	configureLogger()

	go func() {
		logrus.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	config, err := config.LoadConfig()
	if err != nil {
		config.FILEPATH = "./internal/persistence/persisted_messages.txt"
		config.DBPATH = "./sqlite/msgbroker.db"
	}

	file, err := os.OpenFile(config.FILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logrus.Error("Error in opening file: ", err)
	}
	defer file.Close()

	database := initDB(config.DBPATH)
	defer database.Close()

	consumerStorage := channelconsumer.NewInMemoryInMemoryConsumerCache()
	messageQueue := channelconsumer.NewInMemoryMessageQueue()
	consumerIdGenerator := &channelconsumer.SerialConsumerIdGenerator{}
	messageIdGenerator := &channelconsumer.SerialMessageIdGenerator{}
	persist := persistence.New()
	requestCounter := channelconsumer.NewRequestCounter()
	failMsgCounter := channelconsumer.NewFailMsgCounter()
	channel := channelconsumer.NewChannel()
	sqlite := sqlite.New()
	tcpServer := tcpconn.New(":8081", consumerStorage, messageQueue, consumerIdGenerator, messageIdGenerator, persist, file, channel, database, sqlite)

	go func() {
		if err := tcpServer.Listen(); err != nil {
			logrus.Fatalf("Error in starting TCP server: %v", err)
		}
	}()

	app := echo.New()
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	app.POST("/api/channels/:id", func(c echo.Context) error {
		return communication.Broadcast(c, messageQueue, consumerStorage, messageIdGenerator, persist, file, requestCounter, failMsgCounter, channel, database, sqlite)
	})

	app.GET("/api/channel/all", func(c echo.Context) error {
		return table.ChannelDetails(c, messageQueue, consumerStorage, persist, file, requestCounter, failMsgCounter, channel, sqlite, database)
	})

	app.GET("/api/consumer/count", func(c echo.Context) error {
		return chart.Consumer(c, consumerStorage)
	})

	app.GET("/api/message/count", func(c echo.Context) error {
		return chart.Messages(c, messageQueue)
	})

	app.GET("/api/request/count", func(c echo.Context) error {
		return chart.Request(c, requestCounter)
	})

	app.GET("/api/channel/count", func(c echo.Context) error {
		return chart.Channel(c, channel)
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

func initDB(DBPath string) *sql.DB {
	database, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		logrus.Error("Error in opening database: ", err)
	}
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS message (id INTEGER PRIMARY KEY, channel_name TEXT NOT NULL, content BLOB)")
	if err != nil {
		logrus.Error("Error in creating table: ", err)
	}
	return database
}
