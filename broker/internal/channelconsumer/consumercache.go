package channelconsumer

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type ConsumerEvent struct {
	Timestamp time.Time
	Count     int
}

type Storage interface {
	Add(consumer *Consumer)
	Remove(consumerId int, channelName string)
	GetAll() map[string][]Consumer
	Get(consumerId int, channelName string) Consumer
	GetByChannel(channelName string) []Consumer
	GetEventCount() []ConsumerEvent
}

type InMemoryConsumerCache struct {
	mu             sync.RWMutex
	consumers      map[string][]Consumer
	consumerEvents []ConsumerEvent
	sseChannel     chan struct{}
	sseChannSum    chan struct{}
}

func NewInMemoryInMemoryConsumerCache() *InMemoryConsumerCache {
	return &InMemoryConsumerCache{
		consumers:      make(map[string][]Consumer),
		consumerEvents: []ConsumerEvent{},
		sseChannel:     make(chan struct{}, 1),
		sseChannSum:    make(chan struct{}, 1),
	}
}

func (cc *InMemoryConsumerCache) Add(consumer *Consumer) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if cacheConsumers, found := cc.consumers[consumer.SubscribedChannel]; found {
		cacheConsumers := append(cacheConsumers, *consumer)
		cc.consumers[consumer.SubscribedChannel] = cacheConsumers
	} else {
		cc.consumers[consumer.SubscribedChannel] = []Consumer{*consumer}
	}
	cc.recordEvent()
	logrus.Info("Added consumer from cache: ", *consumer)
}

func (cc *InMemoryConsumerCache) Remove(consumerId int, channelName string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cacheConsumers, found := cc.consumers[channelName]
	if !found {
		return
	}

	var updatedConsumers []Consumer
	for _, consumer := range cacheConsumers {
		if consumer.Id != consumerId {
			updatedConsumers = append(updatedConsumers, consumer)
			continue
		}
		logrus.Info("Removed consumerID from cache: ", consumerId)
	}
	cc.consumers[channelName] = updatedConsumers
	if len(updatedConsumers) == 0 {
		delete(cc.consumers, channelName)
	}
	cc.recordEvent()
}

func (cc *InMemoryConsumerCache) GetAll() map[string][]Consumer {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	consumerCacheCopy := make(map[string][]Consumer)
	for channelName, consumer := range cc.consumers {
		consumerCacheCopy[channelName] = append([]Consumer(nil), consumer...)
	}
	return consumerCacheCopy
}

func (cc *InMemoryConsumerCache) Get(consumerId int, channelName string) Consumer {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	consumers, found := cc.consumers[channelName]
	if !found {
		return Consumer{}
	}

	for _, consumer := range consumers {
		if consumer.Id == consumerId {
			return consumer
		}
	}
	return Consumer{}
}

func (cc *InMemoryConsumerCache) GetByChannel(channelName string) []Consumer {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	consumers, found := cc.consumers[channelName]
	if !found {
		return nil
	}
	return append([]Consumer(nil), consumers...)
}

func (cc *InMemoryConsumerCache) recordEvent() {
	count := 0
	for _, consumers := range cc.consumers {
		count += len(consumers)
	}
	cc.consumerEvents = append(cc.consumerEvents, ConsumerEvent{
		Timestamp: time.Now(),
		Count:     count,
	})
	cc.notifySSE()
}

func (cc *InMemoryConsumerCache) GetEventCount() []ConsumerEvent {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	return append([]ConsumerEvent{}, cc.consumerEvents...)
}

func (cc *InMemoryConsumerCache) notifySSE() {
	select {
	case cc.sseChannel <- struct{}{}:
	default:
	}
	select {
	case cc.sseChannSum <- struct{}{}:
	default:
	}
}

func (cc *InMemoryConsumerCache) SSEChannel() <-chan struct{} {
	return cc.sseChannel
}

func (cc *InMemoryConsumerCache) SSEChannelSummary() <-chan struct{} {
	return cc.sseChannSum
}
