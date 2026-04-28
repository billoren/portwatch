package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/notify"
	"portwatch/internal/scanner"
)

func makeEvent(proto, action string) alert.Event {
	p := scanner.Port{Proto: proto, Number: 8080}
	return alert.NewEvent(p, alert.Action(action))
}

func TestWebhookSendsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wn := notify.NewWebhook(notify.WebhookConfig{URL: ts.URL})
	ev := makeEvent("tcp", "denied")
	if err := wn.Notify(ev); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	if received["port"] != "tcp/8080" {
		t.Errorf("expected port tcp/8080, got %v", received["port"])
	}
	if received["action"] != "denied" {
		t.Errorf("expected action denied, got %v", received["action"])
	}
}

func TestWebhookCustomHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.Header.Get("X-Api-Key"); v != "secret" {
			t.Errorf("expected header X-Api-Key=secret, got %q", v)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	wn := notify.NewWebhook(notify.WebhookConfig{
		URL:     ts.URL,
		Headers: map[string]string{"X-Api-Key": "secret"},
	})
	if err := wn.Notify(makeEvent("udp", "new")); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}

func TestWebhookNon2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wn := notify.NewWebhook(notify.WebhookConfig{URL: ts.URL})
	if err := wn.Notify(makeEvent("tcp", "new")); err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestWebhookTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wn := notify.NewWebhook(notify.WebhookConfig{
		URL:     ts.URL,
		Timeout: 50 * time.Millisecond,
	})
	if err := wn.Notify(makeEvent("tcp", "new")); err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}
