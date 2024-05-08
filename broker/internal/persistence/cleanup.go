package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func CleanUp(msg channelconsumer.Message) error {
	filePath := "./internal/persistence/persisted_messages.txt"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %v", err)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error in opening file: %v", err)
	}
	defer file.Close()

	var messages []channelconsumer.Message

	decoder := json.NewDecoder(file)

	for {
		var m channelconsumer.Message
		if err := decoder.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error in decoding JSON: %v", err)
		}
		if m.ChannelId == msg.ChannelId && m.ID == msg.ID {
			continue
		}
		messages = append(messages, m)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking file: %v", err)
	}
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("error truncating file: %v", err)
	}

	encoder := json.NewEncoder(file)

	for _, m := range messages {
		if err := encoder.Encode(m); err != nil {
			return fmt.Errorf("error in encoding JSON: %v", err)
		}
	}

	return nil
}
