package messagequeue

import (
	"fmt"

	"github.com/sagini18/message-broker/internal/channelconsumer"
)

func RemoveMessageFromChannel(msgs []channelconsumer.Message) error {
	channelconsumer.MessageCache.Lock()
	defer channelconsumer.MessageCache.Unlock()

	for _, msg := range msgs {
		channelId := msg.ChannelId
		cachedData, found := channelconsumer.MessageCache.Data[channelId]

		if !found {
			continue
		}

		for index, value := range cachedData {
			if found && value.MessageId == msg.MessageId {
				channelconsumer.MessageCache.Data[channelId] = append(cachedData[:index], cachedData[index+1:]...)

				if len(channelconsumer.MessageCache.Data[channelId]) == 0 {
					delete(channelconsumer.MessageCache.Data, channelId)
				}
				fmt.Println("MessageCache after Deleted: ", channelconsumer.MessageCache.Data)
			}
		}
	}
	return nil
}
