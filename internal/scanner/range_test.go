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
	}

	for _, tc := range tests {
		got, err := scanner.ParsePortList(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParsePortList(%q): expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParsePortList(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if len(got) != len(tc.want) {
			t.Errorf("ParsePortList(%q): got %v, want %v", tc.input, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("ParsePortList(%q)[%d]: got %d, want %d", tc.input, i, got[i], tc.want[i])
			}
		}
	}
}
