package daemon

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestNewDaemonLoadsConfig(t *testing.T) {
	dir := t.TempDir()

	rulesContent := `[{"ports":["80"],"action":"allow"}]`
	rulesFile := writeTempFile(t, dir, "rules.json", rulesContent)

	cfg := map[string]interface{}{
		"ports":        []string{"80"},
		"rules_file":   rulesFile,
		"state_file":   filepath.Join(dir, "state.json"),
		"interval_sec": 30,
		"webhooks":     []interface{}{},
	}
	cfgBytes, _ := json.Marshal(cfg)
	cfgFile := writeTempFile(t, dir, "config.json", string(cfgBytes))

	d, err := New(cfgFile)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil daemon")
	}
	d.ticker.Stop()
}

func TestRunCancels(t *testing.T) {
	dir := t.TempDir()

	rulesContent := `[{"ports":["80"],"action":"allow"}]`
	rulesFile := writeTempFile(t, dir, "rules.json", rulesContent)

	cfg := map[string]interface{}{
		"ports":        []string{"80"},
		"rules_file":   rulesFile,
		"state_file":   filepath.Join(dir, "state.json"),
		"interval_sec": 60,
		"webhooks":     []interface{}{},
	}
	cfgBytes, _ := json.Marshal(cfg)
	cfgFile := writeTempFile(t, dir, "config.json", string(cfgBytes))

	d, err := New(cfgFile)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
