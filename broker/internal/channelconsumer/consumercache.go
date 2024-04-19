package channelconsumer

import (
	"fmt"
	"sync"
)

type Storage interface {
	Add(consumer *Consumer)
	Remove(consumerId int)
	Get() []*Consumer
}

type InMemoryConsumerCache struct {
	mu        sync.Mutex
	consumers []*Consumer //should be a map
}

func NewInMemoryInMemoryConsumerCache() *InMemoryConsumerCache {
	return &InMemoryConsumerCache{
		consumers: make([]*Consumer, 0),
	}
}

func (cc *InMemoryConsumerCache) Add(consumer *Consumer) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.consumers = append(cc.consumers, consumer)

	fmt.Println("--------------------------------------------")
	fmt.Println("ConsumerCache after Added: ", cc.consumers)
}

func (cc *InMemoryConsumerCache) Remove(consumerId int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	for i, c := range cc.consumers {
		if c.Id == consumerId {
			cc.consumers = append(cc.consumers[:i], cc.consumers[i+1:]...)
			fmt.Println("ConsumerCache after Deleted: ", cc.consumers)
			break
		}
	}
}

func (cc *InMemoryConsumerCache) Get() []*Consumer {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	return cc.consumers
}
