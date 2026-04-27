package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a single alert event emitted by the monitor.
type Event struct {
	Timestamp time.Time
	Level     Level
	Port      scanner.Port
	Message   string
}

// String returns a human-readable representation of the event.
func (e Event) String() string {
	return fmt.Sprintf("%s [%s] port %s — %s",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port,
		e.Message,
	)
}

// Handler is a function that receives alert events.
type Handler func(Event)

// Logger returns a Handler that writes events to the supplied writer.
// If w is nil, os.Stderr is used.
func Logger(w io.Writer) Handler {
	if w == nil {
		w = os.Stderr
	}
	return func(e Event) {
		fmt.Fprintln(w, e.String())
	}
}

// Multi returns a Handler that fans an event out to all provided handlers.
func Multi(handlers ...Handler) Handler {
	return func(e Event) {
		for _, h := range handlers {
			h(e)
		}
	}
}

// NewEvent is a convenience constructor.
func NewEvent(level Level, port scanner.Port, msg string) Event {
	return Event{
		Timestamp: time.Now(),
		Level:     level,
		Port:      port,
		Message:   msg,
	}
}
