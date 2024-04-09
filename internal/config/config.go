package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Port             string `json:"port,omitempty"`
	ElecardWebSocket string `json:"elecardWebSocket,omitempty"`
	LogPath          string `json:"logPath,omitempty"`
	TimeoutSec       int    `json:"timeoutSec,omitempty"`
}

func NewConfig(filePath string) Config {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
