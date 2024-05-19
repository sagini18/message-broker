package channelconsumer

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Storage interface {
	Add(consumer *Consumer)
	Remove(consumerId int)
	GetAll() map[int]Consumer
	Get(consumerId int) Consumer
}

type InMemoryConsumerCache struct {
	mu        sync.RWMutex
	consumers map[int]Consumer
}

func NewInMemoryInMemoryConsumerCache() *InMemoryConsumerCache {
	return &InMemoryConsumerCache{
		consumers: make(map[int]Consumer),
	}
}

func (cc *InMemoryConsumerCache) Add(consumer *Consumer) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.consumers[consumer.Id] = *consumer

	logrus.Info("Added consumer from cache: ", *consumer)
}

func (cc *InMemoryConsumerCache) Remove(consumerId int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	delete(cc.consumers, consumerId)
	logrus.Info("Removed consumerID from cache: ", consumerId)
}

func (cc *InMemoryConsumerCache) GetAll() map[int]Consumer {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	// Create a copy of the map to avoid holding the lock during iteration
	consumersCopy := make(map[int]Consumer)
	for id, consumer := range cc.consumers {
		consumersCopy[id] = consumer
	}
	return consumersCopy
}

func (cc *InMemoryConsumerCache) Get(consumerId int) Consumer {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	return cc.consumers[consumerId]
}
