package monitor

import (
	"log"
	"time"

	"portwatch/internal/rules"
	"portwatch/internal/scanner"
)

// Alert represents a detected port change event.
type Alert struct {
	Port   scanner.Port
	Action rules.Action
	Msg    string
}

// Monitor periodically scans ports and emits alerts on unexpected changes.
type Monitor struct {
	scanner  *scanner.Scanner
	ruleSet  *rules.RuleSet
	interval time.Duration
	alerts   chan Alert
	stop     chan struct{}
	prev     map[scanner.Port]bool
}

// New creates a Monitor with the given scanner, rule set, and poll interval.
func New(s *scanner.Scanner, rs *rules.RuleSet, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  s,
		ruleSet:  rs,
		interval: interval,
		alerts:   make(chan Alert, 64),
		stop:     make(chan struct{}),
		prev:     make(map[scanner.Port]bool),
	}
}

// Alerts returns the read-only channel of emitted alerts.
func (m *Monitor) Alerts() <-chan Alert {
	return m.alerts
}

// Start begins the polling loop in a background goroutine.
func (m *Monitor) Start() {
	go m.loop()
}

// Stop signals the polling loop to exit.
func (m *Monitor) Stop() {
	close(m.stop)
}

func (m *Monitor) loop() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := m.tick(); err != nil {
				log.Printf("monitor: scan error: %v", err)
			}
		case <-m.stop:
			return
		}
	}
}

func (m *Monitor) tick() error {
	open, err := m.scanner.Scan()
	if err != nil {
		return err
	}
	current := make(map[scanner.Port]bool, len(open))
	for _, p := range open {
		current[p] = true
		if !m.prev[p] {
			// newly opened port
			action := m.ruleSet.Match(p)
			if action != rules.ActionAllow {
				m.alerts <- Alert{Port: p, Action: action, Msg: "new port opened"}
			}
		}
	}
	for p := range m.prev {
		if !current[p] {
			m.alerts <- Alert{Port: p, Action: rules.ActionAllow, Msg: "port closed"}
		}
	}
	m.prev = current
	return nil
}
