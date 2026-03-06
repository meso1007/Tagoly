package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"tagoly/internal/prompt"
)

type Config struct {
	CustomTags []prompt.CommitType `json:"customTags"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(".tagolycustom"); err == nil {
		data, err := os.ReadFile(".tagolycustom")
		if err != nil {
			if err := json.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("invalid json in .tagolycustom: %w", err)
			}
			return cfg, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, err
	}
	file := filepath.Join(home, ".tagolycustom")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return cfg, nil
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return cfg, err
	}
	json.Unmarshal(data, cfg)
	return cfg, nil
}
