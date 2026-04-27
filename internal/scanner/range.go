package scanner

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePortList parses a comma-separated list of ports and port ranges
// (e.g. "22,80,443,8000-8080") into a deduplicated slice of port numbers.
func ParsePortList(input string) ([]int, error) {
	seen := make(map[int]struct{})
	var ports []int

	for _, part := range strings.Split(input, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid port range %q: %w", part, err)
			}
			hi, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("invalid port range %q: %w", part, err)
			}
			if lo > hi {
				return nil, fmt.Errorf("invalid range %d-%d: low > high", lo, hi)
			}
			for p := lo; p <= hi; p++ {
				if _, ok := seen[p]; !ok {
					seen[p] = struct{}{}
					ports = append(ports, p)
				}
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port %q: %w", part, err)
			}
			if _, ok := seen[p]; !ok {
				seen[p] = struct{}{}
				ports = append(ports, p)
			}
		}
	}
	return ports, nil
}
