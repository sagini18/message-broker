package main

import (
	"database/sql"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagini18/message-broker/broker/config"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/communication"
	"github.com/sagini18/message-broker/broker/internal/tcpconn"
	"github.com/sagini18/message-broker/broker/persistence"
	"github.com/sagini18/message-broker/broker/services"
	"github.com/sirupsen/logrus"
)

func main() {
	configureLogger()

	go func() {
		logrus.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	config, err := config.LoadConfig()
	if err != nil {
		config.DBPATH = "./persistence/msgbroker.db"
	}
	runMigrations(config.DBPATH)

	database := initDB(config.DBPATH)
	defer database.Close()

	consumerStorage := channelconsumer.NewInMemoryInMemoryConsumerCache()
	messageQueue := channelconsumer.NewInMemoryMessageQueue()
	consumerIdGenerator := &channelconsumer.SerialConsumerIdGenerator{}
	messageIdGenerator := &channelconsumer.SerialMessageIdGenerator{}
	requestCounter := channelconsumer.NewRequestCounter()
	failMsgCounter := channelconsumer.NewFailMsgCounter()
	channel := channelconsumer.NewChannel()
	persist := persistence.New()
	tcpServer := tcpconn.New(":8081", consumerStorage, messageQueue, consumerIdGenerator, messageIdGenerator, channel, database, persist)

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

	api := app.Group("/api/v1")

	api.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	api.POST("/channels/:id", func(c echo.Context) error {
		return communication.Broadcast(c, messageQueue, consumerStorage, messageIdGenerator, requestCounter, failMsgCounter, channel, database, persist)
	})

	api.GET("/channels", func(c echo.Context) error {
		return services.Channels(c, messageQueue, consumerStorage, requestCounter, failMsgCounter, channel, persist, database)
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
	return database
}

func runMigrations(DBPath string) {
	m, err := migrate.New(
		"file://migrations",
		"sqlite3://"+DBPath,
	)
	if err != nil {
		logrus.Fatalf("Error creating migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error running migrations: %v", err)
	}
}
