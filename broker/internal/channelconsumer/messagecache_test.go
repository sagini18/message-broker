package channelconsumer

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

func TestMessageAdd(t *testing.T) {
	mockQueue := NewInMemoryMessageQueue()
	mockMessage := Message{
		ID:        1,
		ChannelId: 1,
	}

	mockQueue.Add(mockMessage)

	messages := mockQueue.Get(1)

	assert.Equal(t, 1, len(messages))

	assert.Equal(t, mockMessage, messages[0])
}

func TestMessageRemove(t *testing.T) {
	mockQueue := NewInMemoryMessageQueue()
	mockMessage := Message{
		ID:        1,
		ChannelId: 1,
	}

	mockQueue.Add(mockMessage)

	messages := mockQueue.Get(1)
	assert.Equal(t, 1, len(messages))

	mockQueue.Remove(mockMessage.ID, mockMessage.ChannelId)

	messages = mockQueue.Get(1)
	assert.Equal(t, 0, len(messages))
}

func TestSendPendingMessages(t *testing.T) {
	mockQueue := NewInMemoryMessageQueue()
	mockMessage := Message{
		ID:        1,
		ChannelId: 1,
		Content:   "test",
	}

	mockQueue.Add(mockMessage)

	ConnSpy := &ConnSpy{}

	mockQueue.SendPendingMessages(1, ConnSpy)

	expectedOutput := `[{"ID":1,"ChannelId":1,"Content":"test"}]`

	assert.Equal(t, expectedOutput, ConnSpy.writeBuffer.String())
}

func TestGetAllMessages(t *testing.T) {
	mockQueue := NewInMemoryMessageQueue()
	mockMessage := Message{
		ID:        1,
		ChannelId: 1,
	}

	mockQueue.Add(mockMessage)

	messages := mockQueue.Get(1)

	assert.Equal(t, 1, len(messages))

	assert.Equal(t, mockMessage, messages[0])
}
