package tcpconn

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
)

type TCPServer struct {
	addr                string
	consumerStore       channelconsumer.Storage
	messageQueue        channelconsumer.MessageStorage
	consumerIdGenerator channelconsumer.IdGenerator
	messageIdGenerator  channelconsumer.IdGenerator
	persist             persistence.Persistence
	file                *os.File
}

func New(addr string, store channelconsumer.Storage, msgStore channelconsumer.MessageStorage, consumerIdGenerator channelconsumer.IdGenerator, messageIdGenerator channelconsumer.IdGenerator, persist persistence.Persistence, file *os.File) *TCPServer {
	return &TCPServer{
		addr:                addr,
		consumerStore:       store,
		messageQueue:        msgStore,
		consumerIdGenerator: consumerIdGenerator,
		messageIdGenerator:  messageIdGenerator,
		persist:             persist,
		file:                file,
	}
}

func (t *TCPServer) Listen() error {
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		return fmt.Errorf("tcpserver.Listen() failed: %v", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Errorf("tcpserver.Listen(): listener.Accept error: %v", err)
			continue
		}
		go func() {
			channel, consumer, err := t.handleNewClientConnection(conn)
			if err != nil {
				logrus.Errorf("tcpserver.Listen(): %v", err)
				return
			}

			count := sendPersistedData(channel, conn, t.persist, t.file)
			if count == 0 {
				t.messageQueue.SendPendingMessages(channel, conn)
			}

			if err := listenToConsumerMessages(conn, consumer, t.consumerStore, t.messageQueue, t.persist, t.file); err != nil {
				logrus.Errorf("tcpserver.Listen():  %v", err)
			}
		}()
	}
}

func (t *TCPServer) handleNewClientConnection(connection net.Conn) (int, *channelconsumer.Consumer, error) {
	channelBuf := make([]byte, 200)

	n, err := connection.Read(channelBuf)
	if err != nil {
		return 0, nil, fmt.Errorf("handleNewClientConnection(): reading error from tcp conn: %v", err)
	}

	channel := string(channelBuf[:n])
	channelInt, err := strconv.Atoi(channel)
	if err != nil {
		return 0, nil, fmt.Errorf("handleNewClientConnection(): converting channel into int error: %v", err)
	}

	newId := t.consumerIdGenerator.NewId()
	consumer := channelconsumer.NewConsumer(newId, connection, []int{channelInt})

	t.consumerStore.Add(consumer)

	var channelConfirmation [1]channelconsumer.Message
	channelConfirmation[0] = *channelconsumer.NewMessage(-1, -1, channel)

	confirmationBytes, err := json.Marshal(channelConfirmation)
	if err != nil {
		return 0, nil, fmt.Errorf("handleNewClientConnection(): marshaling error: %v", err)
	}

	if _, err = connection.Write(confirmationBytes); err != nil {
		return 0, nil, fmt.Errorf("handleNewClientConnection(): writing error: %v", err)
	}

	return channelInt, consumer, nil
}

func sendPersistedData(channel int, connection net.Conn, persist persistence.Persistence, file *os.File) int {
	fileData, err := persist.Read(channel, file)
	if err != nil {
		logrus.Errorf("tcpserver.Listen().sendPersistedData(): %v", err)
	}

	if len(fileData) == 0 {
		return 0
	}

	messageBytes, err := json.Marshal(fileData)
	if err != nil {
		logrus.Errorf("tcpserver.Listen().sendPersistedData(): json.Marshal error: %v", err)

	}

	if _, err = connection.Write(messageBytes); err != nil {
		logrus.Errorf("tcpserver.Listen().sendPersistedData(): writing error: %v", err)
	}
	logrus.Info("Sent persisted data to consumer: ", fileData)
	return len(fileData)
}
