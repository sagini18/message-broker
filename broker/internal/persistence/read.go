package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

func Read(channelId int) ([]channelconsumer.Message, error) {
	filePath := "./internal/persistence/persisted_messages.txt"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("Error in opening file: ", err)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var messages []channelconsumer.Message

	for {
		var msg channelconsumer.Message
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			logrus.Error("persistence.Read() : Error in decoding JSON: ", err)
			return nil, err
		}
		if msg.ChannelId != channelId {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func GetLastMessageId() (int, error) {
	filePath := "./internal/persistence/persisted_messages.txt"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return 0, fmt.Errorf("file does not exist: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("Error in opening file: ", err)
		return 0, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	lastId := 0
	for {
		var msg channelconsumer.Message
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			logrus.Error("persistence.GetLastMessageId() : Error in decoding JSON: ", err)
			return 0, err
		}
		if msg.ID > lastId {
			lastId = msg.ID
		}
	}

	return lastId, nil
}
