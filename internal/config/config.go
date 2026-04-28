package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the top-level portwatch daemon configuration.
type Config struct {
	// ScanInterval is how often the port scanner runs.
	ScanInterval Duration `json:"scan_interval"`

	// Ports is the list or range expression of ports to scan.
	// e.g. "22,80,443,8000-9000"
	Ports string `json:"ports"`

	// RulesFile is the path to the JSON rules file.
	RulesFile string `json:"rules_file"`

	// StateFile is the path where port state snapshots are persisted.
	StateFile string `json:"state_file"`

	// LogFile is an optional path to write alert log output.
	// If empty, alerts are written to stdout.
	LogFile string `json:"log_file,omitempty"`
}

// Duration is a time.Duration that marshals/unmarshals as a string (e.g. "30s").
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("config: invalid duration %q: %w", s, err)
	}
	d.Duration = parsed
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

// LoadFile reads a JSON config file from the given path.
func LoadFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()
	return Load(f)
}

// Load decodes a Config from r and validates required fields.
func Load(r interface{ Read([]byte) (int, error) }) (*Config, error) {
	var cfg Config
	if err := json.NewDecoder(r.(interface {
		Read([]byte) (int, error)
	})).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Ports == "" {
		return fmt.Errorf("config: ports must not be empty")
	}
	if c.RulesFile == "" {
		return fmt.Errorf("config: rules_file must not be empty")
	}
	if c.StateFile == "" {
		return fmt.Errorf("config: state_file must not be empty")
	}
	if c.ScanInterval.Duration <= 0 {
		return fmt.Errorf("config: scan_interval must be a positive duration")
	}
	return nil
}
