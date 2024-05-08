package persistence

import (
	"os"

	"github.com/sirupsen/logrus"
)

func AppendToFile(data []byte) error {
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
	return nil
}
