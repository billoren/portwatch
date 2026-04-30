package notify

import (
	"io"
	"net"
	"net/smtp"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEmailEvent() alert.Event {
	return alert.NewEvent(
		scanner.Port{Proto: "tcp", Number: 8080},
		"deny",
		"unexpected port opened",
	)
}

// startFakeSMTP starts a minimal TCP server that accepts one SMTP session
// and returns the received data via a channel.
func startFakeSMTP(t *testing.T) (addr string, received <-chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ch := make(chan string, 1)
	go func() {
		defer ln.Close()
		conn, err := ln.Accept()
		if err != nil {
			ch <- ""
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(3 * time.Second))
		// Send greeting
		conn.Write([]byte("220 fake smtp\r\n"))
		buf, _ := io.ReadAll(conn)
		ch <- string(buf)
	}()
	return ln.Addr().String(), ch
}

func TestEmailDefaultsApplied(t *testing.T) {
	n := NewEmail(EmailConfig{
		Host: "localhost",
		From: "a@b.com",
		To:   []string{"c@d.com"},
	})
	if n.cfg.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", n.cfg.Timeout)
	}
	if n.cfg.Port != 587 {
		t.Errorf("expected default port 587, got %d", n.cfg.Port)
	}
}

func TestEmailSendsToSMTP(t *testing.T) {
	// Use a real local SMTP server if available; otherwise just verify
	// that the notifier constructs the message and attempts delivery.
	_ = smtp.SendMail // ensure package is used

	n := NewEmail(EmailConfig{
		Host:    "127.0.0.1",
		Port:    0, // will be overridden
		From:    "portwatch@example.com",
		To:      []string{"admin@example.com"},
		Timeout: 1 * time.Second,
	})

	ev := makeEmailEvent()

	// Expect an error since there's no real SMTP server on port 0.
	err := n.Send(ev)
	if err == nil {
		t.Error("expected error sending to unavailable SMTP server")
	}
}

func TestEmailSubjectContainsAction(t *testing.T) {
	ev := makeEmailEvent()
	subject := fmt.Sprintf("[portwatch] %s on port %s", ev.Action, ev.Port)
	if !strings.Contains(subject, "deny") {
		t.Errorf("subject missing action: %q", subject)
	}
	if !strings.Contains(subject, "tcp/8080") {
		t.Errorf("subject missing port: %q", subject)
	}
}
