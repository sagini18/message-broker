package channelconsumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsumerAdd(t *testing.T) {
	mockStorage := NewInMemoryInMemoryConsumerCache()
	addedConsumer := &Consumer{
		Id:                1,
		SubscribedChannel: "test",
	}

	mockStorage.Add(addedConsumer)

	consumers := mockStorage.GetAll()
	assert.Equal(t, 1, len(consumers))
	assert.Equal(t, *addedConsumer, consumers["test"][0])
}

func TestConsumerRemove(t *testing.T) {
	mockStorage := NewInMemoryInMemoryConsumerCache()
	mockConsumer := &Consumer{
		Id:                1,
		SubscribedChannel: "test",
	}

	mockStorage.Add(mockConsumer)

	consumers := mockStorage.GetAll()
	assert.Equal(t, 1, len(consumers))

	mockStorage.Remove(1, "test")

	consumers = mockStorage.GetAll()
	assert.Equal(t, 0, len(consumers))
}

func TestGetConsumer(t *testing.T) {
	mockStorage := NewInMemoryInMemoryConsumerCache()
	mockConsumer := &Consumer{
		Id:                1,
		SubscribedChannel: "test",
	}

	mockStorage.Add(mockConsumer)

	var consumer Consumer = mockStorage.Get(1, "test")
	assert.Equal(t, *mockConsumer, consumer)
}

func TestGet(t *testing.T) {
	mockStorage := NewInMemoryInMemoryConsumerCache()
	mockConsumer := &Consumer{
		Id:                1,
		SubscribedChannel: "test",
	}

	mockStorage.Add(mockConsumer)

	consumers := mockStorage.GetAll()
	assert.Equal(t, 1, len(consumers))
	assert.Equal(t, *mockConsumer, consumers["test"][0])
}
