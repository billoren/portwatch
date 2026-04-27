package scanner_test

import (
	"net"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/scanner"
)

func TestScanDetectsOpenPort(t *testing.T) {
	// Start a temporary TCP listener on an ephemeral port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	s := scanner.New([]int{port}, 500*time.Millisecond)

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	found := false
	for _, p := range ports {
		if p.Port == port && p.Protocol == "tcp" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected port %d/tcp to be detected as open", port)
	}
}

func TestScanClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed and requires no privilege to check.
	s := scanner.New([]int{1}, 200*time.Millisecond)
	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected no open ports, got %v", ports)
	}
}

func TestPortString(t *testing.T) {
	p := scanner.Port{Protocol: "tcp", Address: "127.0.0.1", Port: 8080}
	expected := "tcp://127.0.0.1:8080"
	if p.String() != expected {
		t.Errorf("expected %q, got %q", expected, p.String())
	}
}
