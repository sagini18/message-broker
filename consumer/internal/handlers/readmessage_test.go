package handlers

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ConnSpy struct {
	writeBuffer bytes.Buffer
}

func (c *ConnSpy) Read(b []byte) (n int, err error) {
	copy(b, []byte(`[{"MessageId": 1,"ChannelId": 1,"Content": "Hello World"}]`))
	trimmedSlice := bytes.TrimRightFunc(b, func(r rune) bool {
		return r == 0
	})
	return len(trimmedSlice), nil
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

type ConnSpyMock struct {
	Conn ConnSpy
}

func (c *ConnSpyMock) New(conn ConnSpy) {
	c.Conn = conn
}

func TestReadAndUnmarshalMessage(t *testing.T) {
	connSpy := ConnSpy{}
	connSpyMock := &ConnSpyMock{}
	connSpyMock.New(connSpy)

	buffer, totalBytesRead, err := readAndExpandBuffer(&connSpyMock.Conn)
	assert.Nil(t, err)

	b := make([]byte, 200)
	copy(b, []byte(`[{"MessageId": 1,"ChannelId": 1,"Content": "Hello World"}]`))
	trimmedSlice := bytes.TrimRightFunc(b, func(r rune) bool {
		return r == 0
	})
	assert.Equal(t, len(trimmedSlice), totalBytesRead)
	assert.Equal(t, b, buffer)
}
