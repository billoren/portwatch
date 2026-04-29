package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the top-level portwatch daemon configuration.
type Config struct {
	Ports     string        `json:"ports"`
	RulesFile string        `json:"rules_file"`
	Interval  time.Duration `json:"interval"`
	StateFile string        `json:"state_file"`
	Webhooks  []WebhookCfg  `json:"webhooks,omitempty"`
	Slack     *SlackCfg     `json:"slack,omitempty"`
}

// WebhookCfg configures a generic HTTP webhook notifier.
type WebhookCfg struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Timeout time.Duration     `json:"timeout"`
}

// SlackCfg configures the Slack incoming-webhook notifier.
type SlackCfg struct {
	WebhookURL string        `json:"webhook_url"`
	Timeout    time.Duration `json:"timeout"`
}

// LoadFile reads and parses a JSON config file from path.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}
	return Load(data)
}

// Load parses a JSON-encoded config from raw bytes.
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
	if cfg.Ports == "" {
		return fmt.Errorf("config: 'ports' is required")
	}
	if cfg.RulesFile == "" {
		return fmt.Errorf("config: 'rules_file' is required")
	}
	if cfg.Interval <= 0 {
		return fmt.Errorf("config: 'interval' must be positive")
	}
	if cfg.Slack != nil && cfg.Slack.WebhookURL == "" {
		return fmt.Errorf("config: slack 'webhook_url' is required when slack block is present")
	}
	return nil
}
