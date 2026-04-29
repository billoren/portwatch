package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"portwatch/internal/alert"
)

// SlackNotifier sends alert events to a Slack incoming webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlack creates a SlackNotifier that posts to the given Slack webhook URL.
func NewSlack(webhookURL string, timeout time.Duration) *SlackNotifier {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: timeout},
	}
}

// Send formats the event as a Slack message and posts it to the webhook URL.
func (s *SlackNotifier) Send(ctx context.Context, ev alert.Event) error {
	payload := slackPayload{
		Text: fmt.Sprintf(":warning: *portwatch alert* — %s", ev.String()),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("slack: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
