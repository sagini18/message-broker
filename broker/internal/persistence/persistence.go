package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

type Persistence interface {
	Write(data []byte, file *os.File) error
	Read(channelId string, file *os.File) ([]channelconsumer.Message, error)
	ReadAll(file *os.File) (map[string][]channelconsumer.Message, error)
	Remove(messageId int, file *os.File) error
	ReadCount(channelName string, file *os.File) int
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
		return fmt.Errorf("persistence.Write(): %v", err)
	}
	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("persistence.Write() newline error: %v", err)
	}

	return nil
}

func (p *persistence) Read(channelId string, file *os.File) ([]channelconsumer.Message, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("persistence.Read() : seeking file error: %v", err)
	}

	decoder := json.NewDecoder(file)

	var messages []channelconsumer.Message

	for {
		var msg channelconsumer.Message
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("persistence.Read() : error decoding JSON: %v", err)
		}
		if msg.ChannelName != channelId {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (p *persistence) ReadAll(file *os.File) (map[string][]channelconsumer.Message, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		logrus.Error("Error in seeking file: ", err)
		return nil, err
	}

	decoder := json.NewDecoder(file)

	messages := make(map[string][]channelconsumer.Message)

	for {
		var msg channelconsumer.Message
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			logrus.Error("persistence.ReadAll() : Error in decoding JSON: ", err)
			return nil, err
		}
		if _, found := messages[msg.ChannelName]; !found {
			messages[msg.ChannelName] = []channelconsumer.Message{msg}
			continue
		}
		messages[msg.ChannelName] = append(messages[msg.ChannelName], msg)
	}

	return messages, nil
}

func (p *persistence) Remove(messageID int, file *os.File) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var modifiedContent []channelconsumer.Message

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("persistence.Remove() : seeking file error: %v", err)
	}

	decoder := json.NewDecoder(file)

	for {
		var msg channelconsumer.Message
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("persistence.Remove() : error decoding JSON: %v", err)
		}
		if msg.ID != messageID {
			modifiedContent = append(modifiedContent, msg)
		}
	}

	newFile, err := os.OpenFile(file.Name(), os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("persistence.Remove() : error opening file: %v", err)
	}

	encoder := json.NewEncoder(newFile)

	for _, msg := range modifiedContent {
		if err := encoder.Encode(msg); err != nil {
			return fmt.Errorf("persistence.Remove() : error encoding JSON: %v", err)
		}
	}

	return nil
}

func (p *persistence) ReadCount(channelName string, file *os.File) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		logrus.Error("Error in seeking file: ", err)
		return 0
	}

	decoder := json.NewDecoder(file)

	count := 0
	for {
		var msg channelconsumer.Message
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			logrus.Error("persist.ReadCount() : error decoding JSON: ", err)
			return 0
		}
		if msg.ChannelName == channelName {
			count++
		}
	}

	return count
}
