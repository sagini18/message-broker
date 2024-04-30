package main

import (
	"os"

	"github.com/sagini18/message-broker/consumer/internal/handlers"
	"github.com/sagini18/message-broker/consumer/internal/types"

	"github.com/sirupsen/logrus"
)

func main() {
	configureLogger()

	tcpConn := &types.TcpConn{}
	receiver := &types.Receiver{}

	if err := handlers.InitConnection(tcpConn); err != nil {
		logrus.Fatal("Error in connecting to the server: handlers.InitConnection(): ", err)
	}

	if err := handlers.HandleChannel(tcpConn.Conn); err != nil {
		logrus.Error("Error in handling channel: handlers.HandleChannel(): ", err)
		return
	}

	handlers.ReadMessage(tcpConn.Conn, receiver)
}

func configureLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
