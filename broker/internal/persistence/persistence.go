package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

type Persistence interface {
	Write(data []byte) error
	Read(channelId int) ([]channelconsumer.Message, error)
	CleanUp(messageId int, channelId int) error
}

type persistence struct {
}

func New() Persistence {
	return &persistence{}
}

func (p *persistence) Write(data []byte) error {
	filePath := "./internal/persistence/persisted_messages.txt"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Error("Error in opening file: ", err)
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		logrus.Error("Error in writing to file: ", err)
		return err
	}
	if _, err := file.WriteString("\n"); err != nil {
		logrus.Error("Error in writing newline: ", err)
		return err
	}

	return nil
}

func (p *persistence) Read(channelId int) ([]channelconsumer.Message, error) {
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

func (p *persistence) CleanUp(messageId int, channelId int) error {
	filePath := "./internal/persistence/persisted_messages.txt"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %v", err)
	}

	inputFile, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		// if err.Error() == "open "+filePath+": The process cannot access the file because it is being used by another process." {
		// 	time.Sleep(10 * time.Second)
		// 	inputFile, err = os.OpenFile(filePath, os.O_RDONLY, 0644)
		// }
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer inputFile.Close()

	tempFilePath := filePath + ".temp"
	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating temporary output file: %v", err)
	}
	defer tempFile.Close()

	decoder := json.NewDecoder(inputFile)

	encoder := json.NewEncoder(tempFile)

	for decoder.More() {
		var m channelconsumer.Message
		if err := decoder.Decode(&m); err != nil {
			return fmt.Errorf("error decoding JSON: %v", err)
		}

		if m.ChannelId == channelId && m.ID == messageId {
			continue
		}

		if err := encoder.Encode(m); err != nil {
			return fmt.Errorf("error encoding JSON: %v", err)
		}
	}

	if err := inputFile.Close(); err != nil {
		return fmt.Errorf("error closing input file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("error closing temporary output file: %v", err)
	}

	if err := os.Rename(tempFilePath, filePath); err != nil {
		if err.Error() == "rename "+tempFilePath+" "+filePath+": Access is denied." {
			time.Sleep(10 * time.Second)
			if err := os.Rename(tempFilePath, filePath); err == nil {
				return nil
			}
		}
		return fmt.Errorf("error renaming temporary file: %v", err)
	}
	return nil
}
