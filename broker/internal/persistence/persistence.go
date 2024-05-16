package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/sagini18/message-broker/broker/config"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

type Persistence interface {
	Write(data []byte) error
	Read(channelId int) ([]channelconsumer.Message, error)
	Remove(messageId int) error
}

type persistence struct {
	filePath string
	mu       sync.Mutex
}

func New() Persistence {
	config, err := config.LoadConfig()
	if err != nil {
		config.FilePath = "./internal/persistence/persisted_messages.txt"
	}
	return &persistence{
		filePath: config.FilePath,
	}
}

func (p *persistence) Write(data []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	file, err := os.OpenFile(p.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := os.Stat(p.filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %v", err)
	}

	file, err := os.OpenFile(p.filePath, os.O_RDONLY, 0644)
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

func (p *persistence) Remove(messageID int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	file, err := os.OpenFile(p.filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("persist.Remove() : error opening file : %v", err)
	}
	defer file.Close()

	var modifiedContent []channelconsumer.Message

	decoder := json.NewDecoder(file)

	for {
		var msg channelconsumer.Message

		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("persist.Remove() : error decoding JSON: %v", err)
		}

		if msg.ID == messageID {
			continue
		}

		modifiedContent = append(modifiedContent, msg)
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("persist.Remove() : error truncating file: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("persist.Remove() : error seeking to file beginning: %v", err)
	}

	encoder := json.NewEncoder(file)

	for _, msg := range modifiedContent {
		if err := encoder.Encode(msg); err != nil {
			return fmt.Errorf("error encoding JSON: %v", err)
		}
	}

	return nil
}
