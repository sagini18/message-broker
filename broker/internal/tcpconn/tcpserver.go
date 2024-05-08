package tcpconn

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

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
}

func New(addr string, store channelconsumer.Storage, msgStore channelconsumer.MessageStorage, consumerIdGenerator channelconsumer.IdGenerator, messageIdGenerator channelconsumer.IdGenerator) *TCPServer {
	return &TCPServer{
		addr:                addr,
		consumerStore:       store,
		messageQueue:        msgStore,
		consumerIdGenerator: consumerIdGenerator,
		messageIdGenerator:  messageIdGenerator,
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
				logrus.Errorf("tcpserver.Listen(): handleNewClientConnection failed to %v: %v", conn.RemoteAddr().String(), err)
				return
			}

			count := sendPersistedData(channel, conn)
			if count == 0 {
				t.messageQueue.SendPendingMessages(channel, conn)
			}

			if err := listenToConsumerMessages(conn, consumer, t.consumerStore, t.messageQueue); err != nil {
				logrus.Errorf("tcpserver.Listen(): listenToConsumerMessages failed to %v: %v", conn.RemoteAddr().String(), err)
			}
		}()
	}
}

func (t *TCPServer) handleNewClientConnection(connection net.Conn) (int, *channelconsumer.Consumer, error) {
	channelBuf := make([]byte, 200)

	n, err := connection.Read(channelBuf)
	if err != nil {
		return 0, nil, fmt.Errorf("handleNewClientConnection: reading error from tcp conn: %v", err)
	}

	channel := string(channelBuf[:n])
	channelInt, err := strconv.Atoi(channel)
	if err != nil {
		return 0, nil, fmt.Errorf("handleNewClientConnection: converting channel into int error: %v", err)
	}

	newId := t.consumerIdGenerator.NewId()
	consumer := channelconsumer.NewConsumer(newId, connection, []int{channelInt})

	t.consumerStore.Add(consumer)

	var channelConfirmation [1]channelconsumer.Message
	channelConfirmation[0] = *channelconsumer.NewMessage(-1, -1, channel)

	confirmationBytes, err := json.Marshal(channelConfirmation)
	if err != nil {
		return 0, nil, fmt.Errorf("handleNewClientConnection: marshaling error: %v", err)
	}

	if _, err = connection.Write(confirmationBytes); err != nil {
		return 0, nil, fmt.Errorf("handleNewClientConnection: writing error: %v", err)
	}

	return channelInt, consumer, nil
}

func sendPersistedData(channel int, connection net.Conn) int {
	time.Sleep(10 * time.Second)
	fileData, err := persistence.Read(channel)
	if err != nil {
		logrus.Errorf("tcpserver.Listen(): persistence.Read() failed: %v", err)
	}

	if len(fileData) == 0 {
		return 0
	}

	messageBytes, err := json.Marshal(fileData)
	if err != nil {
		logrus.Errorf("sendPersistedData(): json.Marshal error: %v", err)

	}

	if _, err = connection.Write(messageBytes); err != nil {
		logrus.Errorf("sendPersistedData(): writing error: %v", err)
	}
	logrus.Info("Sent persisted data to consumer: ", fileData)
	return len(fileData)
}
