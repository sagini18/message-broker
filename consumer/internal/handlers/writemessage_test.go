package handlers

import (
	"testing"

	"github.com/sagini18/message-broker/consumer/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestWriteMessage(t *testing.T) {
	connSpy := ConnSpy{}
	connSpyMock := &ConnSpyMock{}
	connSpyMock.New(connSpy)

	receiver := types.Receiver{
		ReceivedMessage: []byte(`[{"MessageId": 1,"ChannelId": 1,"Content": "Hello World"}]`),
	}

	WriteMessage(&connSpyMock.Conn, &receiver)

	assert.Equal(t, `[{"MessageId": 1,"ChannelId": 1,"Content": "Hello World"}]`, connSpyMock.Conn.writeBuffer.String())
}
