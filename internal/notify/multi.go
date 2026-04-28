package notify

import (
	"errors"
	"fmt"

	"portwatch/internal/alert"
)

// Notifier is the interface for sending alert events.
type Notifier interface {
	Notify(ev alert.Event) error
}

// Multi fans an event out to multiple Notifiers, collecting all errors.
type Multi struct {
	notifiers []Notifier
}

// NewMulti returns a Multi that dispatches to each provided Notifier.
func NewMulti(nn ...Notifier) *Multi {
	return &Multi{notifiers: nn}
}

// Add appends a Notifier to the fan-out list.
func (m *Multi) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Notify calls every registered Notifier and returns a joined error if any fail.
func (m *Multi) Notify(ev alert.Event) error {
	var errs []error
	for i, n := range m.notifiers {
		if err := n.Notify(ev); err != nil {
			errs = append(errs, fmt.Errorf("notifier[%d]: %w", i, err))
		}
	}
	return errors.Join(errs...)
}
