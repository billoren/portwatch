package rules

import (
	"testing"
)

func baseRuleSet() *RuleSet {
	return &RuleSet{
		DefaultAction: ActionAlert,
		Rules: []Rule{
			{Name: "allow-web", Ports: []int{80, 443}, Action: ActionAllow},
			{Name: "alert-db", Ports: []int{5432, 3306}, Action: ActionAlert},
			{Name: "deny-telnet", Ports: []int{23}, Action: ActionDeny},
		},
	}
}

func TestMatchKnownPort(t *testing.T) {
	rs := baseRuleSet()
	action, name := rs.Match(80)
	if action != ActionAllow || name != "allow-web" {
		t.Errorf("expected allow/allow-web, got %s/%s", action, name)
	}
}

func TestMatchDenyPort(t *testing.T) {
	rs := baseRuleSet()
	action, name := rs.Match(23)
	if action != ActionDeny || name != "deny-telnet" {
		t.Errorf("expected deny/deny-telnet, got %s/%s", action, name)
	}
}

func TestMatchDefaultFallback(t *testing.T) {
	rs := baseRuleSet()
	action, name := rs.Match(9999)
	if action != ActionAlert || name != "default" {
		t.Errorf("expected alert/default, got %s/%s", action, name)
	}
}

// TestMatchSecondPort verifies that a rule matches on any port in its list,
// not just the first one.
func TestMatchSecondPort(t *testing.T) {
	rs := baseRuleSet()
	action, name := rs.Match(443)
	if action != ActionAllow || name != "allow-web" {
		t.Errorf("expected allow/allow-web for port 443, got %s/%s", action, name)
	}
}

func TestValidateValid(t *testing.T) {
	rs := baseRuleSet()
	if err := rs.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateInvalidDefaultAction(t *testing.T) {
	rs := baseRuleSet()
	rs.DefaultAction = "unknown"
	if err := rs.Validate(); err == nil {
		t.Fatal("expected validation error for invalid default action")
	}
}

func TestValidateMissingPorts(t *testing.T) {
	rs := &RuleSet{
		DefaultAction: ActionAlert,
		Rules: []Rule{
			{Name: "empty", Ports: []int{}, Action: ActionAllow},
		},
	}
	if err := rs.Validate(); err == nil {
		t.Fatal("expected validation error for rule with no ports")
	}
}
