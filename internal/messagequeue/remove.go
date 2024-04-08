package messagequeue

import (
	"fmt"
	"strconv"

	"github.com/sagini18/message-broker/internal/message"
)

func RemoveMessageFromChannel(msgs []message.Message) error {
	message.MessageCache.Lock()
	defer message.MessageCache.Unlock()

	for _, msg := range msgs {
		channelId := strconv.Itoa(msg.ChannelId)
		cachedData, found := message.MessageCache.Data[channelId]

		if !found {
			continue
		}

		for index, value := range cachedData {
			if found && value.MessageId == msg.MessageId {
				message.MessageCache.Data[channelId] = append(cachedData[:index], cachedData[index+1:]...)

				if len(message.MessageCache.Data[channelId]) == 0 {
					delete(message.MessageCache.Data, channelId)
				}
				fmt.Println("MessageCache after Deleted: ", message.MessageCache.Data)
			}
		}
	}
	return nil
}
