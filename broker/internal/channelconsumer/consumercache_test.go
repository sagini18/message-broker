package channelconsumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsumerAdd(t *testing.T) {
	mockStorage := NewInMemoryInMemoryConsumerCache()
	addedConsumer := &Consumer{
		Id:                 1,
		SubscribedChannels: []int{1},
	}

	mockStorage.Add(addedConsumer)

	consumers := mockStorage.Get()
	assert.Equal(t, 1, len(consumers))
	assert.Equal(t, *addedConsumer, consumers[1])
}

func TestConsumerRemove(t *testing.T) {
	mockStorage := NewInMemoryInMemoryConsumerCache()
	mockConsumer := &Consumer{
		Id:                 1,
		SubscribedChannels: []int{1},
	}

	mockStorage.Add(mockConsumer)

	consumers := mockStorage.Get()
	assert.Equal(t, 1, len(consumers))

	mockStorage.Remove(1)

	consumers = mockStorage.Get()
	assert.Equal(t, 0, len(consumers))
}

func TestGetConsumer(t *testing.T) {
	mockStorage := NewInMemoryInMemoryConsumerCache()
	mockConsumer := &Consumer{
		Id:                 1,
		SubscribedChannels: []int{1},
	}

	mockStorage.Add(mockConsumer)

	consumer := mockStorage.GetConsumer(1)
	assert.Equal(t, *mockConsumer, consumer)
}

func TestGet(t *testing.T) {
	mockStorage := NewInMemoryInMemoryConsumerCache()
	mockConsumer := &Consumer{
		Id:                 1,
		SubscribedChannels: []int{1},
	}

	mockStorage.Add(mockConsumer)

	consumers := mockStorage.Get()
	assert.Equal(t, 1, len(consumers))
	assert.Equal(t, *mockConsumer, consumers[1])
}
