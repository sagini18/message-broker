package messagequeue

import (
	"fmt"
	"strconv"

	"github.com/sagini18/message-broker/internal/message"
)

func RemoveMessageFromChannel(msgs []message.Message) error {
	MessageCache.Lock()
	defer MessageCache.Unlock()

	for _, msg := range msgs {
		channelId := strconv.Itoa(msg.ChannelId)
		cachedData, found := MessageCache.Data[channelId]

		if !found {
			return fmt.Errorf("CHANNEL_NOT_FOUND")
		}

		for index, value := range cachedData {
			if found && value.MessageId == msg.MessageId {
				MessageCache.Data[channelId] = append(cachedData[:index], cachedData[index+1:]...)

				if len(MessageCache.Data[channelId]) == 0 {
					delete(MessageCache.Data, channelId)
				}
				fmt.Println("MessageCache after Deleted: ", MessageCache.Data)
				fmt.Println("-------------------------------------------------------------------")
			}
		}
	}
	return nil
}
