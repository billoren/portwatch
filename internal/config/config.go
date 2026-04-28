package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// WebhookConfig holds settings for a single outbound webhook.
type WebhookConfig struct {
	URL        string            `json:"url"`
	Headers    map[string]string `json:"headers"`
	TimeoutSec int               `json:"timeout_sec"`
}

// Config is the top-level daemon configuration.
type Config struct {
	Ports       []string        `json:"ports"`
	RulesFile   string          `json:"rules_file"`
	StateFile   string          `json:"state_file"`
	IntervalSec int             `json:"interval_sec"`
	Webhooks    []WebhookConfig `json:"webhooks"`
}

// LoadFile reads and parses a JSON config from path.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}
	return Load(data)
}

// Load parses a JSON config from raw bytes.
func Load(data []byte) (*Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	if len(cfg.Ports) == 0 {
		return errors.New("config: ports must not be empty")
	}
	if cfg.RulesFile == "" {
		return errors.New("config: rules_file must not be empty")
	}
	if cfg.IntervalSec <= 0 {
		return errors.New("config: interval_sec must be > 0")
	}
	for i, wh := range cfg.Webhooks {
		if wh.URL == "" {
			return fmt.Errorf("config: webhooks[%d]: url must not be empty", i)
		}
		if wh.TimeoutSec <= 0 {
			cfg.Webhooks[i].TimeoutSec = 5
		}
	}
	if cfg.StateFile == "" {
		cfg.StateFile = "portwatch-state.json"
	}
	return nil
}
