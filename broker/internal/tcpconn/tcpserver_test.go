package tcpconn

import (
	"net"
	"testing"
	"time"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/stretchr/testify/assert"
)

type MockConn struct{}

func (m *MockConn) Read(b []byte) (n int, err error) {
	copy(b, []byte("123"))
	return len("123"), nil
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (m *MockConn) Close() error {
	return nil
}

func (m *MockConn) LocalAddr() net.Addr {
	return nil
}

func (m *MockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestHandleNewClientConnection(t *testing.T) {
	mockStorage := channelconsumer.NewInMemoryInMemoryConsumerCache()
	mockMessageQueue := channelconsumer.NewInMemoryMessageQueue()
	mockConsumerIdGenerator := &channelconsumer.SerialConsumerIdGenerator{}
	mockMessageIdGenerator := &channelconsumer.SerialMessageIdGenerator{}

	server := New(":8081", mockStorage, mockMessageQueue, mockConsumerIdGenerator, mockMessageIdGenerator)

	mockConn := &MockConn{}

	channel, consumer, err := server.handleNewClientConnection(mockConn)
	assert.NoError(t, err)

	expectedChannel := 123
	assert.Equal(t, expectedChannel, channel)

	expectedConsumerId := 1
	if assert.NotNil(t, consumer) {
		assert.Equal(t, expectedConsumerId, consumer.Id)
		assert.Equal(t, expectedChannel, consumer.SubscribedChannels[0])
	}
}