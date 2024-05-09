package channelconsumer

import (
	"sync/atomic"

	"github.com/google/uuid"
)

type IdGenerator interface {
	NewId() int
}

type SerialConsumerIdGenerator struct {
	lastId atomic.Uint32
}

func (s *SerialConsumerIdGenerator) NewId() int {
	id := s.lastId.Add(1)
	return int(id)
}

type SerialMessageIdGenerator struct {
}

func (s *SerialMessageIdGenerator) NewId() int {
	newUUID := uuid.New().ID()
	return int(newUUID)
}
