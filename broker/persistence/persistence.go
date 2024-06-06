package persistence

import (
	"database/sql"
	"fmt"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
	"github.com/sirupsen/logrus"
)

type Persistence interface {
	Write(data channelconsumer.Message, db *sql.DB) error
	Read(channelName string, db *sql.DB) ([]channelconsumer.Message, error)
	ReadAll(db *sql.DB) (map[string][]channelconsumer.Message, error)
	Remove(messageId int, db *sql.DB) error
	ReadCount(channelName string, db *sql.DB) int
}

type persistence struct {
	mu sync.Mutex
}

func New() Persistence {
	return &persistence{}
}

func (p *persistence) Write(data channelconsumer.Message, db *sql.DB) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	p.mu.Lock()
	defer p.mu.Unlock()

	content, err := json.Marshal(data.Content)
	if err != nil {
		return fmt.Errorf("persistence.Write() marshalling error: %v", err)
	}

	query := `INSERT INTO message (id,channel_name, content) VALUES (?,?, ?)`
	_, err = db.Exec(query, data.ID, data.ChannelName, content)
	if err != nil {
		return fmt.Errorf("persistence.Write() inserting into table error: %v", err)
	}

	return nil
}

func (p *persistence) Read(channelName string, db *sql.DB) ([]channelconsumer.Message, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	p.mu.Lock()
	defer p.mu.Unlock()

	query := `SELECT id, channel_name, content FROM message WHERE channel_name = ?`
	rows, err := db.Query(query, channelName)
	if err != nil {
		return nil, fmt.Errorf("persistence.Read() query error: %v", err)
	}
	defer rows.Close()

	var messages []channelconsumer.Message

	for rows.Next() {
		var id int
		var channelName string
		var content []byte
		if err := rows.Scan(&id, &channelName, &content); err != nil {
			return nil, fmt.Errorf("persistence.Read() scan error: %v", err)
		}

		var msgContent interface{}
		if err := json.Unmarshal(content, &msgContent); err != nil {
			return nil, fmt.Errorf("persistence.Read() unmarshal content error: %v", err)
		}

		msg := channelconsumer.Message{
			ID:          id,
			ChannelName: channelName,
			Content:     msgContent,
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (p *persistence) ReadAll(db *sql.DB) (map[string][]channelconsumer.Message, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	p.mu.Lock()
	defer p.mu.Unlock()

	query := `SELECT id, channel_name, content FROM message`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("persistence.Read() query error: %v", err)
	}
	defer rows.Close()

	messages := make(map[string][]channelconsumer.Message)

	for rows.Next() {
		var id int
		var channelName string
		var content []byte
		if err := rows.Scan(&id, &channelName, &content); err != nil {
			return nil, fmt.Errorf("persistence.Read() scan error: %v", err)
		}

		var msgContent interface{}
		if err := json.Unmarshal(content, &msgContent); err != nil {
			return nil, fmt.Errorf("persistence.Read() unmarshal content error: %v", err)
		}

		msg := channelconsumer.Message{
			ID:          id,
			ChannelName: channelName,
			Content:     msgContent,
		}
		messages[msg.ChannelName] = append(messages[msg.ChannelName], msg)
	}
	return messages, nil
}

func (p *persistence) Remove(messageID int, db *sql.DB) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	query := `DELETE FROM message WHERE id = ?`
	_, err := db.Exec(query, messageID)
	if err != nil {
		return fmt.Errorf("persistence.Remove() error: %v", err)
	}

	return nil
}

func (p *persistence) ReadCount(channelName string, db *sql.DB) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	query := `SELECT COUNT(*) FROM message WHERE channel_name = ?`
	var count int
	err := db.QueryRow(query, channelName).Scan(&count)
	if err != nil {
		logrus.Error("persistence.ReadCount() error: ", err)
		return 0
	}

	return count
}
