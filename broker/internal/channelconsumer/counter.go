package channelconsumer

import (
	"sync"
	"sync/atomic"
)

type Counter interface {
	Add(channelName string)
	Get(channelName string) int
}

type RequestCounter struct {
	mu    sync.RWMutex
	count map[string]*atomic.Uint32
}

func NewRequestCounter() *RequestCounter {
	return &RequestCounter{
		count: make(map[string]*atomic.Uint32),
	}
}

func (p *RequestCounter) Add(channelName string) {
	p.mu.Lock()
	if _, exists := p.count[channelName]; !exists {
		p.count[channelName] = new(atomic.Uint32)
	}
	p.count[channelName].Add(1)
	p.mu.Unlock()
}

func (p *RequestCounter) Get(channelName string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if count, exists := p.count[channelName]; exists {
		return int(count.Load())
	}
	return 0
}

type FailMsgCounter struct {
	mu    sync.RWMutex
	count map[string]*atomic.Uint32
}

func NewFailMsgCounter() *FailMsgCounter {
	return &FailMsgCounter{
		count: make(map[string]*atomic.Uint32),
	}
}

func (f *FailMsgCounter) Add(channelName string) {
	f.mu.Lock()
	if _, exists := f.count[channelName]; !exists {
		f.count[channelName] = new(atomic.Uint32)
	}
	f.count[channelName].Add(1)
	f.mu.Unlock()
}

func (f *FailMsgCounter) Get(channelName string) int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if count, exists := f.count[channelName]; exists {
		return int(count.Load())
	}
	return 0
}

func (f *FailMsgCounter) GetAllChannel() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	var channels []string
	for channel := range f.count {
		channels = append(channels, channel)
	}
	return channels
}
