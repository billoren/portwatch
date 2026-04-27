package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port with its protocol and address.
type Port struct {
	Protocol string
	Address  string
	Port     int
}

// String returns a human-readable representation of the port.
func (p Port) String() string {
	return fmt.Sprintf("%s://%s:%d", p.Protocol, p.Address, p.Port)
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	Timeout time.Duration
	Ports   []int
}

// New creates a new Scanner with the given port range and timeout.
func New(ports []int, timeout time.Duration) *Scanner {
	return &Scanner{
		Timeout: timeout,
		Ports:   ports,
	}
}

// Scan checks which ports in the configured list are open.
// It returns a slice of open Port entries.
func (s *Scanner) Scan() ([]Port, error) {
	var open []Port
	for _, p := range s.Ports {
		for _, proto := range []string{"tcp", "udp"} {
			addr := fmt.Sprintf("127.0.0.1:%d", p)
			conn, err := net.DialTimeout(proto, addr, s.Timeout)
			if err == nil {
				conn.Close()
				open = append(open, Port{
					Protocol: proto,
					Address:  "127.0.0.1",
					Port:     p,
				})
			}
		}
	}
	return open, nil
}
