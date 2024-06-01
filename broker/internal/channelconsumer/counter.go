package channelconsumer

import (
	"sync"
	"sync/atomic"
	"time"
)

type Counter interface {
	Add(channelName string)
	Get(channelName string) int
}

type RequestEvent struct {
	Timestamp time.Time
	Count     uint32
}

type RequestCounter struct {
	mu            sync.RWMutex
	count         map[string]*atomic.Uint32
	requestEvents []RequestEvent
	sseChannel    chan struct{}
	sseChannSum   chan struct{}
}

func NewRequestCounter() *RequestCounter {
	return &RequestCounter{
		count:         make(map[string]*atomic.Uint32),
		requestEvents: []RequestEvent{},
		sseChannel:    make(chan struct{}, 1),
		sseChannSum:   make(chan struct{}, 1),
	}
}

func (rc *RequestCounter) Add(channelName string) {
	rc.mu.Lock()
	if _, exists := rc.count[channelName]; !exists {
		rc.count[channelName] = new(atomic.Uint32)
	}
	rc.count[channelName].Add(1)
	rc.mu.Unlock()

	rc.recordEvent()
}

func (rc *RequestCounter) Get(channelName string) int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if count, exists := rc.count[channelName]; exists {
		return int(count.Load())
	}
	return 0
}

func (rc *RequestCounter) recordEvent() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if len(rc.requestEvents) < 1 {
		rc.requestEvents = append(rc.requestEvents, RequestEvent{
			Timestamp: time.Now(),
			Count:     1,
		})
		rc.notifySSE()
		return
	}
	lastNo := rc.requestEvents[len(rc.requestEvents)-1].Count

	rc.requestEvents = append(rc.requestEvents, RequestEvent{
		Timestamp: time.Now(),
		Count:     lastNo + 1,
	})

	rc.notifySSE()
}

func (rc *RequestCounter) GetEventCount() []RequestEvent {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return append([]RequestEvent{}, rc.requestEvents...)
}

func (rc *RequestCounter) SSEChannel() <-chan struct{} {
	return rc.sseChannel
}

func (rc *RequestCounter) SSEChannelSummary() <-chan struct{} {
	return rc.sseChannSum
}

func (rc *RequestCounter) notifySSE() {
	select {
	case rc.sseChannel <- struct{}{}:
	default:
	}
	select {
	case rc.sseChannSum <- struct{}{}:
	default:
	}
}

type FailMsgCounter struct {
	mu          sync.RWMutex
	count       map[string]*atomic.Uint32
	sseChannel  chan struct{}
	sseChannSum chan struct{}
}

func NewFailMsgCounter() *FailMsgCounter {
	return &FailMsgCounter{
		count:       make(map[string]*atomic.Uint32),
		sseChannel:  make(chan struct{}, 1),
		sseChannSum: make(chan struct{}, 1),
	}
}

func (f *FailMsgCounter) Add(channelName string) {
	f.mu.Lock()
	if _, exists := f.count[channelName]; !exists {
		f.count[channelName] = new(atomic.Uint32)
	}
	f.count[channelName].Add(1)
	f.mu.Unlock()

	f.notifySSE()
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

func (f *FailMsgCounter) SSEChannel() <-chan struct{} {
	return f.sseChannel
}

func (f *FailMsgCounter) SSEChannelSummary() <-chan struct{} {
	return f.sseChannSum
}

func (f *FailMsgCounter) notifySSE() {
	select {
	case f.sseChannel <- struct{}{}:
	default:
	}
	select {
	case f.sseChannSum <- struct{}{}:
	default:
	}
}
