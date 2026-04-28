package rules

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/portwatch/internal/scanner"
)

// rawRule is the JSON representation of a rule.
type rawRule struct {
	Name    string `json:"name"`
	Ports   string `json:"ports"`
	Action  string `json:"action"`
	Comment string `json:"comment"`
}

// rawConfig is the top-level JSON config structure.
type rawConfig struct {
	DefaultAction string    `json:"default_action"`
	Rules         []rawRule `json:"rules"`
}

// LoadFile reads a JSON rules file and returns a validated RuleSet.
func LoadFile(path string) (*RuleSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading rules file: %w", err)
	}
	return Load(data)
}

// Load parses JSON bytes into a validated RuleSet.
func Load(data []byte) (*RuleSet, error) {
	var cfg rawConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing rules JSON: %w", err)
	}

	rs := &RuleSet{
		DefaultAction: Action(cfg.DefaultAction),
	}

	for i, raw := range cfg.Rules {
		if raw.Name == "" {
			return nil, fmt.Errorf("rule at index %d: name must not be empty", i)
		}
		ports, err := scanner.ParsePortList(raw.Ports)
		if err != nil {
			return nil, fmt.Errorf("rule %q: invalid ports %q: %w", raw.Name, raw.Ports, err)
		}
		rs.Rules = append(rs.Rules, Rule{
			Name:    raw.Name,
			Ports:   ports,
			Action:  Action(raw.Action),
			Comment: raw.Comment,
		})
	}

	if err := rs.Validate(); err != nil {
		return nil, fmt.Errorf("invalid ruleset: %w", err)
	}
	return rs, nil
}
