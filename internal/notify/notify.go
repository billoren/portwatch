package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"portwatch/internal/alert"
)

// WebhookConfig holds configuration for a webhook notifier.
type WebhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Timeout time.Duration     `json:"timeout,omitempty"`
}

// WebhookNotifier sends alert events to a remote HTTP endpoint.
type WebhookNotifier struct {
	cfg    WebhookConfig
	client *http.Client
}

// NewWebhook creates a new WebhookNotifier with the given config.
func NewWebhook(cfg WebhookConfig) *WebhookNotifier {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &WebhookNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}
}

type payload struct {
	Timestamp time.Time `json:"timestamp"`
	Port      string    `json:"port"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
}

// Notify sends the event to the configured webhook URL.
func (w *WebhookNotifier) Notify(ev alert.Event) error {
	p := payload{
		Timestamp: ev.Time,
		Port:      ev.Port.String(),
		Action:    string(ev.Action),
		Message:   ev.String(),
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, w.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.cfg.Headers {
		req.Header.Set(k, v)
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("notify: http post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
