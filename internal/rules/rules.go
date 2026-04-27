package rules

import (
	"fmt"
	"strings"
)

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAllow Action = "allow"
	ActionAlert Action = "alert"
	ActionDeny  Action = "deny"
)

// Rule represents a single port monitoring rule.
type Rule struct {
	Name    string
	Ports   []int
	Action  Action
	Comment string
}

// RuleSet holds a collection of rules and a default action.
type RuleSet struct {
	Rules         []Rule
	DefaultAction Action
}

// Match returns the action and rule name for a given port.
// Falls back to DefaultAction if no rule matches.
func (rs *RuleSet) Match(port int) (Action, string) {
	for _, r := range rs.Rules {
		for _, p := range r.Ports {
			if p == port {
				return r.Action, r.Name
			}
		}
	}
	return rs.DefaultAction, "default"
}

// Validate checks that all rules have valid actions and at least one port.
func (rs *RuleSet) Validate() error {
	valid := map[Action]bool{ActionAllow: true, ActionAlert: true, ActionDeny: true}
	if !valid[rs.DefaultAction] {
		return fmt.Errorf("invalid default action: %q", rs.DefaultAction)
	}
	for _, r := range rs.Rules {
		if strings.TrimSpace(r.Name) == "" {
			return fmt.Errorf("rule missing name")
		}
		if !valid[r.Action] {
			return fmt.Errorf("rule %q has invalid action: %q", r.Name, r.Action)
		}
		if len(r.Ports) == 0 {
			return fmt.Errorf("rule %q has no ports", r.Name)
		}
	}
	return nil
}
