package communication

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/stretchr/testify/assert"
)

type MockConn struct {
	writeBuffer bytes.Buffer
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	copy(b, []byte("123"))
	return len("123"), nil
}

func (c *MockConn) Write(b []byte) (int, error) {
	return c.writeBuffer.Write(b)
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

func TestBroadcast(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/api/channels/123", strings.NewReader(`{"content": "Hello, World!"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	mockConn := &MockConn{}

	consumerStorage := channelconsumer.NewInMemoryInMemoryConsumerCache()
	consumerStorage.Add(&channelconsumer.Consumer{Id: 1, SubscribedChannels: []int{123}, TcpConn: mockConn})
	messageIdGenerator := &channelconsumer.SerialMessageIdGenerator{}
	messageQueue := channelconsumer.NewInMemoryMessageQueue()

	err := Broadcast(c, messageQueue, consumerStorage, messageIdGenerator)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	allMessags := messageQueue.Get()
	assert.Equal(t, 1, len(allMessags[123]))
	assert.Equal(t, "Hello, World!", allMessags[123][0].Content)
}
