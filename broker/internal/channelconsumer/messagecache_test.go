package channelconsumer

import (
	"bytes"
	"net"
	"testing"
	"time"

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

func TestMessageAdd(t *testing.T) {
	mockQueue := NewInMemoryMessageQueue()
	mockMessage := Message{
		ID:        1,
		ChannelId: 1,
	}

	mockQueue.Add(mockMessage)

	messages := mockQueue.Get()

	assert.Equal(t, 1, len(messages))

	assert.Equal(t, mockMessage, messages[1][0])
}

func TestMessageRemove(t *testing.T) {
	mockQueue := NewInMemoryMessageQueue()
	mockMessage := Message{
		ID:        1,
		ChannelId: 1,
	}

	mockQueue.Add(mockMessage)

	messages := mockQueue.Get()
	assert.Equal(t, 1, len(messages))

	mockQueue.Remove(mockMessage)

	messages = mockQueue.Get()
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

	mockConn := &MockConn{}

	mockQueue.SendPendingMessages(1, mockConn)

	expectedOutput := `[{"ID":1,"ChannelId":1,"Content":"test"}]`

	assert.Equal(t, expectedOutput, mockConn.writeBuffer.String())
}

func TestGetAllMessages(t *testing.T) {
	mockQueue := NewInMemoryMessageQueue()
	mockMessage := Message{
		ID:        1,
		ChannelId: 1,
	}

	mockQueue.Add(mockMessage)

	messages := mockQueue.Get()

	assert.Equal(t, 1, len(messages))

	assert.Equal(t, mockMessage, messages[1][0])
}
