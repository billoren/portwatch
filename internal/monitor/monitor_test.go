package monitor_test

import (
	"testing"
	"time"

	"portwatch/internal/monitor"
	"portwatch/internal/rules"
	"portwatch/internal/scanner"
)

// stubScanner returns a fixed list of open ports on each call.
type stubScanner struct {
	results [][]scanner.Port
	call    int
}

func (s *stubScanner) Scan() ([]scanner.Port, error) {
	if s.call >= len(s.results) {
		return s.results[len(s.results)-1], nil
	}
	res := s.results[s.call]
	s.call++
	return res, nil
}

func makePort(n uint16) scanner.Port { return scanner.Port{Number: n, Proto: "tcp"} }

func TestAlertOnNewDeniedPort(t *testing.T) {
	rs := &rules.RuleSet{
		Rules: []rules.Rule{
			{Ports: []uint16{9999}, Action: rules.ActionDeny},
		},
		Default: rules.ActionAllow,
	}
	stub := &stubScanner{
		results: [][]scanner.Port{
			{},                    // first tick: nothing open
			{makePort(9999)},      // second tick: denied port appears
		},
	}
	sc := scanner.WrapStub(stub)
	m := monitor.New(sc, rs, 10*time.Millisecond)
	m.Start()
	defer m.Stop()

	select {
	case a := <-m.Alerts():
		if a.Port.Number != 9999 {
			t.Fatalf("expected alert on port 9999, got %d", a.Port.Number)
		}
		if a.Action != rules.ActionDeny {
			t.Fatalf("expected deny action, got %v", a.Action)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for alert")
	}
}

func TestAlertOnPortClosed(t *testing.T) {
	rs := &rules.RuleSet{Default: rules.ActionAllow}
	stub := &stubScanner{
		results: [][]scanner.Port{
			{makePort(8080)},
			{},
		},
	}
	sc := scanner.WrapStub(stub)
	m := monitor.New(sc, rs, 10*time.Millisecond)
	m.Start()
	defer m.Stop()

	select {
	case a := <-m.Alerts():
		if a.Port.Number != 8080 {
			t.Fatalf("expected alert on port 8080, got %d", a.Port.Number)
		}
		if a.Msg != "port closed" {
			t.Fatalf("unexpected message: %s", a.Msg)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for closed-port alert")
	}
}
