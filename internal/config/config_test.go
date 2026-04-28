package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValid(t *testing.T) {
	data := `{"ports":["80","443"],"rules_file":"rules.json","state_file":"state.json","interval_sec":10}`
	cfg, err := Load([]byte(data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(cfg.Ports))
	}
	if cfg.IntervalSec != 10 {
		t.Errorf("expected interval 10, got %d", cfg.IntervalSec)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	_, err := Load([]byte(`{bad json}`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadMissingPorts(t *testing.T) {
	data := `{"rules_file":"rules.json","interval_sec":10}`
	_, err := Load([]byte(data))
	if err == nil {
		t.Fatal("expected error for missing ports")
	}
}

func TestLoadMissingRulesFile(t *testing.T) {
	data := `{"ports":["80"],"interval_sec":10}`
	_, err := Load([]byte(data))
	if err == nil {
		t.Fatal("expected error for missing rules_file")
	}
}

func TestLoadZeroInterval(t *testing.T) {
	data := `{"ports":["80"],"rules_file":"rules.json","interval_sec":0}`
	_, err := Load([]byte(data))
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestDefaultStateFile(t *testing.T) {
	data := `{"ports":["80"],"rules_file":"rules.json","interval_sec":5}`
	cfg, err := Load([]byte(data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StateFile == "" {
		t.Error("expected default state_file to be set")
	}
}

func TestWebhookDefaultTimeout(t *testing.T) {
	data := `{"ports":["80"],"rules_file":"r.json","interval_sec":5,"webhooks":[{"url":"http://example.com","timeout_sec":0}]}`
	cfg, err := Load([]byte(data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Webhooks[0].TimeoutSec != 5 {
		t.Errorf("expected default timeout 5, got %d", cfg.Webhooks[0].TimeoutSec)
	}
}

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.json")
	content := `{"ports":["22"],"rules_file":"r.json","interval_sec":15}`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFile(p)
	if err != nil {
		t.Fatalf("LoadFile error: %v", err)
	}
	if cfg.Ports[0] != "22" {
		t.Errorf("expected port 22, got %s", cfg.Ports[0])
	}
}
