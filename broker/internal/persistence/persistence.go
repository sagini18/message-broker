package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gofrs/flock"
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
	lock     *flock.Flock
	mu       sync.Mutex
}

func New() Persistence {
	config, err := config.LoadConfig()
	if err != nil {
		config.FilePath = "./internal/persistence/persisted_messages.txt"
	}
	return &persistence{
		filePath: config.FilePath,
		lock:     flock.New(config.FilePath),
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

	// if err := p.lock.Lock(); err != nil {
	// 	return fmt.Errorf("error in acquiring lock persistence.Write(): %v", err)
	// }
	// defer p.lock.Unlock()

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

	// if err := p.lock.RLock(); err != nil {
	// 	return nil, fmt.Errorf("error in acquiring lock persistence.Read(): %v", err)
	// }
	// defer p.lock.Unlock()

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
func (p *persistence) Remove(messageId int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	tempFilePath := p.filePath + ".temp"

	inputFile, err := os.Open(p.filePath)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer inputFile.Close()

	tempFile, err := os.Create(tempFilePath)
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
		if m.ID != messageId {
			if err := encoder.Encode(m); err != nil {
				return fmt.Errorf("error encoding JSON: %v", err)
			}
		}
	}

	if err := inputFile.Close(); err != nil {
		return fmt.Errorf("error closing input file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("error closing temporary file: %v", err)
	}

	if err := os.Rename(tempFilePath, p.filePath); err != nil {
		return fmt.Errorf("error renaming temporary file: %v", err)
	}

	return nil
}
