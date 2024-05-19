package channelconsumer

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Storage interface {
	Add(consumer *Consumer)
	Remove(consumerId int, channelName string)
	GetAll() map[string][]Consumer
	Get(consumerId int, channelName string) Consumer
	GetByChannel(channelName string) []Consumer
}

type InMemoryConsumerCache struct {
	mu        sync.RWMutex
	consumers map[string][]Consumer
}

func NewInMemoryInMemoryConsumerCache() *InMemoryConsumerCache {
	return &InMemoryConsumerCache{
		consumers: make(map[string][]Consumer),
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
}

func (cc *InMemoryConsumerCache) GetAll() map[string][]Consumer {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	return cc.consumers
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

	return cc.consumers[channelName]
}
