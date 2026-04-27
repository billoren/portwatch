package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, number uint16) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestEventString(t *testing.T) {
	e := alert.Event{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Level:     alert.LevelAlert,
		Port:      makePort("tcp", 8080),
		Message:   "unexpected open port",
	}
	got := e.String()
	if !strings.Contains(got, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", got)
	}
	if !strings.Contains(got, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", got)
	}
	if !strings.Contains(got, "unexpected open port") {
		t.Errorf("expected message in output, got: %s", got)
	}
}

func TestLoggerWritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	h := alert.Logger(&buf)
	e := alert.NewEvent(alert.LevelWarn, makePort("udp", 53), "dns port opened")
	h(e)

	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("expected WARN in logger output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "dns port opened") {
		t.Errorf("expected message in logger output, got: %s", buf.String())
	}
}

func TestMultiFanout(t *testing.T) {
	called := make([]bool, 3)
	handlers := make([]alert.Handler, 3)
	for i := range handlers {
		idx := i
		handlers[idx] = func(_ alert.Event) { called[idx] = true }
	}

	multi := alert.Multi(handlers...)
	multi(alert.NewEvent(alert.LevelInfo, makePort("tcp", 22), "ssh detected"))

	for i, c := range called {
		if !c {
			t.Errorf("handler %d was not called", i)
		}
	}
}

func TestNewEventTimestamp(t *testing.T) {
	before := time.Now()
	e := alert.NewEvent(alert.LevelInfo, makePort("tcp", 80), "test")
	after := time.Now()

	if e.Timestamp.Before(before) || e.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", e.Timestamp, before, after)
	}
}
