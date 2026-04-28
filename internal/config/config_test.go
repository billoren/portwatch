package config

import (
	"strings"
	"testing"
	"time"
)

const validJSON = `{
  "scan_interval": "30s",
  "ports": "22,80,443,8000-9000",
  "rules_file": "/etc/portwatch/rules.json",
  "state_file": "/var/lib/portwatch/state.json"
}`

func TestLoadValid(t *testing.T) {
	cfg, err := Load(strings.NewReader(validJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval.Duration != 30*time.Second {
		t.Errorf("scan_interval: got %v, want 30s", cfg.ScanInterval.Duration)
	}
	if cfg.Ports != "22,80,443,8000-9000" {
		t.Errorf("ports: got %q", cfg.Ports)
	}
	if cfg.RulesFile != "/etc/portwatch/rules.json" {
		t.Errorf("rules_file: got %q", cfg.RulesFile)
	}
	if cfg.StateFile != "/var/lib/portwatch/state.json" {
		t.Errorf("state_file: got %q", cfg.StateFile)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	_, err := Load(strings.NewReader(`{bad json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadMissingPorts(t *testing.T) {
	raw := `{"scan_interval":"10s","rules_file":"r.json","state_file":"s.json"}`
	_, err := Load(strings.NewReader(raw))
	if err == nil {
		t.Fatal("expected error when ports is empty")
	}
}

func TestLoadMissingRulesFile(t *testing.T) {
	raw := `{"scan_interval":"10s","ports":"80","state_file":"s.json"}`
	_, err := Load(strings.NewReader(raw))
	if err == nil {
		t.Fatal("expected error when rules_file is empty")
	}
}

func TestLoadZeroInterval(t *testing.T) {
	raw := `{"scan_interval":"0s","ports":"80","rules_file":"r.json","state_file":"s.json"}`
	_, err := Load(strings.NewReader(raw))
	if err == nil {
		t.Fatal("expected error for zero scan_interval")
	}
}

func TestLoadInvalidDuration(t *testing.T) {
	raw := `{"scan_interval":"notaduration","ports":"80","rules_file":"r.json","state_file":"s.json"}`
	_, err := Load(strings.NewReader(raw))
	if err == nil {
		t.Fatal("expected error for invalid duration string")
	}
}

func TestDurationRoundTrip(t *testing.T) {
	cfg, err := Load(strings.NewReader(validJSON))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	b, err := cfg.ScanInterval.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var d Duration
	if err := d.UnmarshalJSON(b); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if d.Duration != 30*time.Second {
		t.Errorf("round-trip: got %v, want 30s", d.Duration)
	}
}
