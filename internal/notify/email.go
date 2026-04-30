package notify

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// EmailConfig holds the SMTP configuration for email notifications.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
	Timeout  time.Duration
}

// emailNotifier sends alert events via SMTP email.
type emailNotifier struct {
	cfg EmailConfig
}

// NewEmail creates a new email notifier with the given configuration.
func NewEmail(cfg EmailConfig) *emailNotifier {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	return &emailNotifier{cfg: cfg}
}

// Send delivers the alert event as an email to all configured recipients.
func (e *emailNotifier) Send(ev alert.Event) error {
	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)

	subject := fmt.Sprintf("[portwatch] %s on port %s", ev.Action, ev.Port)
	body := fmt.Sprintf(
		"Time: %s\nPort: %s\nAction: %s\nMessage: %s\n",
		ev.Time.Format(time.RFC3339),
		ev.Port,
		ev.Action,
		ev.Message,
	)

	msg := strings.Join([]string{
		"From: " + e.cfg.From,
		"To: " + strings.Join(e.cfg.To, ", "),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
		body,
	}, "\r\n")

	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)
	}

	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, []byte(msg))
}
