package persistence

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

type Persistence interface {
	Write(data []byte, file *os.File) error
	Read(channelId string, file *os.File) ([]channelconsumer.Message, error)
	Remove(messageId int, file *os.File) error
}

type persistence struct {
	mu sync.Mutex
}

func New() Persistence {
	return &persistence{}
}

func (p *persistence) Write(data []byte, file *os.File) error {
	p.mu.Lock()
	defer p.mu.Unlock()

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

func (p *persistence) Read(channelId string, file *os.File) ([]channelconsumer.Message, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		logrus.Error("Error in seeking file: ", err)
		return nil, err
	}

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
		if msg.ChannelName != channelId {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (p *persistence) Remove(messageID int, file *os.File) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var modifiedContent []channelconsumer.Message

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		logrus.Error("Error in seeking file: ", err)
		return err
	}

	decoder := json.NewDecoder(file)

	for {
		var msg channelconsumer.Message
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			logrus.Error("persist.Remove() : error decoding JSON: ", err)
			return err
		}
		if msg.ID != messageID {
			modifiedContent = append(modifiedContent, msg)
		}
	}

	newFile, err := os.OpenFile(file.Name(), os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		logrus.Error("Error opening file for truncation: ", err)
		return err
	}

	encoder := json.NewEncoder(newFile)

	for _, msg := range modifiedContent {
		if err := encoder.Encode(msg); err != nil {
			logrus.Error("persist.Remove() : error encoding JSON: ", err)
			return err
		}
	}

	return nil
}
