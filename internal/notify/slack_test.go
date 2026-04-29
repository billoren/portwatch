package notify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/scanner"
)

func makeSlackEvent() alert.Event {
	return alert.NewEvent("denied", scanner.Port{Proto: "tcp", Number: 9090})
}

func TestSlackSendsJSON(t *testing.T) {
	var received slackPayload
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

	sn := NewSlack(ts.URL, 3*time.Second)
	ev := makeSlackEvent()
	if err := sn.Send(context.Background(), ev); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if !strings.Contains(received.Text, "portwatch alert") {
		t.Errorf("expected alert text, got: %s", received.Text)
	}
	if !strings.Contains(received.Text, "tcp/9090") {
		t.Errorf("expected port in text, got: %s", received.Text)
	}
}

func TestSlackNon2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	sn := NewSlack(ts.URL, 3*time.Second)
	if err := sn.Send(context.Background(), makeSlackEvent()); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSlackTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	sn := NewSlack(ts.URL, 50*time.Millisecond)
	if err := sn.Send(context.Background(), makeSlackEvent()); err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestSlackDefaultTimeout(t *testing.T) {
	// NewSlack with zero timeout should not panic and should set a default.
	sn := NewSlack("http://localhost", 0)
	if sn.client.Timeout <= 0 {
		t.Error("expected positive default timeout")
	}
}
