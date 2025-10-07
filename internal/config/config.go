package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type LoggingConfig struct {
	Level      string `json:"level"`
	OutputFile string `json:"output_file"`
	MaxSizeMB  int64  `json:"max_size_mb"`
	Console    bool   `json:"console"`
}

type Config struct {
	Logging LoggingConfig `json:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	config, err := loadConfigFromFile(configPath)
	if err != nil {
		return nil, err
	}
	return config, nil

}

func loadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "INFO"
	}
	if config.Logging.OutputFile == "" {
		config.Logging.OutputFile = "code-tools-mcp.log"
	}
	if config.Logging.MaxSizeMB == 0 {
		config.Logging.MaxSizeMB = 10
	}

	if config.Logging.OutputFile != "" {
		config.Logging.Console = true
	}

	return &config, nil
}
