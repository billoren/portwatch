package scanner_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/scanner"
)

func TestParsePortList(t *testing.T) {
	tests := []struct {
		input    string
		want     []int
		wantErr  bool
	}{
		{"80", []int{80}, false},
		{"22,80,443", []int{22, 80, 443}, false},
		{"8000-8003", []int{8000, 8001, 8002, 8003}, false},
		{"22,8000-8002,443", []int{22, 8000, 8001, 8002, 443}, false},
		{"80,80", []int{80}, false}, // deduplication
		{"", []int(nil), false},
		{"abc", nil, true},
		{"8000-abc", nil, true},
		{"9000-8000", nil, true}, // lo > hi
		{"0", nil, true},         // port 0 is invalid
		{"65535", []int{65535}, false},
		{"65536", nil, true}, // exceeds max port
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := scanner.ParsePortList(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("ParsePortList(%q): expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("ParsePortList(%q): unexpected error: %v", tc.input, err)
				return
			}
			if len(got) != len(tc.want) {
				t.Errorf("ParsePortList(%q): got %v, want %v", tc.input, got, tc.want)
				return
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("ParsePortList(%q)[%d]: got %d, want %d", tc.input, i, got[i], tc.want[i])
				}
			}
		})
	}
}
