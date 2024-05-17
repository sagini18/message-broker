package communication

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/config"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sagini18/message-broker/broker/internal/persistence"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type ConnSpy struct {
	writeBuffer bytes.Buffer
}

func (c *ConnSpy) Read(b []byte) (n int, err error) {
	copy(b, []byte("123"))
	return len("123"), nil
}

func (c *ConnSpy) Write(b []byte) (int, error) {
	return c.writeBuffer.Write(b)
}
func (c *ConnSpy) Close() error {
	return nil
}

func (c *ConnSpy) LocalAddr() net.Addr {
	return nil
}

func (c *ConnSpy) RemoteAddr() net.Addr {
	return nil
}

func (c *ConnSpy) SetDeadline(t time.Time) error {
	return nil
}

func (c *ConnSpy) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *ConnSpy) SetWriteDeadline(t time.Time) error {
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

	ConnSpy := &ConnSpy{}

	consumerStorage := channelconsumer.NewInMemoryInMemoryConsumerCache()
	consumerStorage.Add(&channelconsumer.Consumer{Id: 1, SubscribedChannels: []int{123}, TcpConn: ConnSpy})
	messageIdGenerator := &channelconsumer.SerialMessageIdGenerator{}
	messageQueue := channelconsumer.NewInMemoryMessageQueue()
	persist := persistence.New()

	config, err := config.LoadConfig()
	if err != nil {
		config.FilePath = "./internal/persistence/persisted_messages.txt"
	}
	file, err := os.OpenFile(config.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Error("Error in opening file: ", err)
	}
	defer file.Close()

	err = Broadcast(c, messageQueue, consumerStorage, messageIdGenerator, persist, file)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, "[{\"ID\":1,\"ChannelId\":123,\"Content\":\"Hello, World!\"}]\n", rec.Body.String())

	allMessags := messageQueue.Get(123)
	assert.Equal(t, 1, len(allMessags))
	assert.Equal(t, "Hello, World!", allMessags[0].Content)
}

func BenchmarkBroadcast(b *testing.B) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/api/channels/123", strings.NewReader(`{"content": "Hello, World!"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	connSpy := &ConnSpy{}

	consumerStorage := channelconsumer.NewInMemoryInMemoryConsumerCache()
	consumerStorage.Add(&channelconsumer.Consumer{Id: 1, SubscribedChannels: []int{123}, TcpConn: connSpy})
	messageIdGenerator := &channelconsumer.SerialMessageIdGenerator{}
	messageQueue := channelconsumer.NewInMemoryMessageQueue()
	persist := persistence.New()

	config, err := config.LoadConfig()
	if err != nil {
		config.FilePath = "./internal/persistence/persisted_messages.txt"
	}
	file, err := os.OpenFile(config.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Error("Error in opening file: ", err)
	}
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := Broadcast(c, messageQueue, consumerStorage, messageIdGenerator, persist, file)
		assert.Nil(b, err)
	}
}
