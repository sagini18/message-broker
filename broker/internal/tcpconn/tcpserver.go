package tcpconn

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

type TCPServer struct {
	addr                string
	consumerStore       channelconsumer.Storage
	messageQueue        channelconsumer.MessageQueue
	consumerIdGenerator channelconsumer.IdGenerator
	messageIdGenerator  channelconsumer.IdGenerator
}

func New(addr string, store channelconsumer.Storage, queue channelconsumer.MessageQueue, consumerIdGenerator channelconsumer.IdGenerator, messageIdGenerator channelconsumer.IdGenerator) *TCPServer {
	return &TCPServer{
		addr:                addr,
		consumerStore:       store,
		messageQueue:        queue,
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

			t.messageQueue.SendPendingMessages(channel, conn)

			go func() {
				if err := listenToConsumerMessages(conn, consumer, t.consumerStore); err != nil {
					logrus.Errorf("tcpserver.Listen(): listenToConsumerMessages failed to %v: %v", conn.RemoteAddr().String(), err)
				}
			}()

			removeMessages(t.messageQueue)
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
