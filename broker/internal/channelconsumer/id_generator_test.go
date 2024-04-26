package channelconsumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsumerId(t *testing.T) {
	var s SerialConsumerIdGenerator
	id := s.NewId()
	assert.Equal(t, 1, id)

	secondId := s.NewId()
	assert.Equal(t, 2, secondId)
}

func TestNewMessageId(t *testing.T) {
	var s SerialMessageIdGenerator
	id := s.NewId()
	assert.Equal(t, 1, id)

	secondId := s.NewId()
	assert.Equal(t, 2, secondId)
}
