package channelconsumer

import (
	"sync/atomic"
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
	lastId atomic.Uint32
}

func (s *SerialMessageIdGenerator) NewId() int {
	id := s.lastId.Add(1)
	return int(id)
}

func (s *SerialMessageIdGenerator) SetLastId(id int) {
	s.lastId.Store(uint32(id))
}
