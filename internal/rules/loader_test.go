package rules

import (
	"testing"
)

var validJSON = []byte(`{
  "default_action": "alert",
  "rules": [
    {"name": "allow-web",   "ports": "80,443",       "action": "allow", "comment": "HTTP/S"},
    {"name": "deny-telnet", "ports": "23",           "action": "deny"},
    {"name": "alert-range", "ports": "8000-8005",    "action": "alert"}
  ]
}`)

func TestLoadValid(t *testing.T) {
	rs, err := Load(validJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs.Rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rs.Rules))
	}
	if rs.DefaultAction != ActionAlert {
		t.Errorf("expected default alert, got %s", rs.DefaultAction)
	}
}

func TestLoadRangeExpansion(t *testing.T) {
	rs, err := Load(validJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// alert-range covers 8000-8005 (6 ports)
	var alertRange *Rule
	for i := range rs.Rules {
		if rs.Rules[i].Name == "alert-range" {
			alertRange = &rs.Rules[i]
		}
	}
	if alertRange == nil {
		t.Fatal("alert-range rule not found")
	}
	if len(alertRange.Ports) != 6 {
		t.Errorf("expected 6 ports in range, got %d", len(alertRange.Ports))
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	_, err := Load([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadInvalidAction(t *testing.T) {
	bad := []byte(`{"default_action":"allow","rules":[{"name":"x","ports":"22","action":"bogus"}]}`)
	_, err := Load(bad)
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestLoadInvalidPorts(t *testing.T) {
	bad := []byte(`{"default_action":"allow","rules":[{"name":"x","ports":"abc","action":"allow"}]}`)
	_, err := Load(bad)
	if err == nil {
		t.Fatal("expected error for invalid port spec")
	}
}
