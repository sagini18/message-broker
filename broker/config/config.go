package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	FilePath string
}

func LoadConfig() (Config, error) {
	var config Config

	file, err := os.OpenFile("./config/config.json", os.O_RDONLY, 0644)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
