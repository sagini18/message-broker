package channelconsumer

type MessageQueue interface {
	Add(message Message)
	Remove(queue MessageStorage)
}

type MessageChannel struct {
	MessageChann chan Message
}

func NewMessageChannel() *MessageChannel {
	return &MessageChannel{
		MessageChann: make(chan Message),
	}
}

func (mc *MessageChannel) Add(msg Message) {
	mc.MessageChann <- msg
}

func (mc *MessageChannel) Remove(store MessageStorage) {
	for msg := range mc.MessageChann {
		store.Remove(msg)
	}
}
