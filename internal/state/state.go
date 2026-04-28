package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a recorded set of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time      `json:"timestamp"`
	Ports     []scanner.Port `json:"ports"`
}

// SaveSnapshot writes the current port list to a JSON file at path.
func SaveSnapshot(path string, ports []scanner.Port) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// LoadSnapshot reads a previously saved snapshot from path.
// Returns an empty Snapshot (with nil Ports) if the file does not exist.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Diff compares two port slices and returns ports that were added or removed.
func Diff(prev, curr []scanner.Port) (added, removed []scanner.Port) {
	prevSet := make(map[scanner.Port]struct{}, len(prev))
	for _, p := range prev {
		prevSet[p] = struct{}{}
	}
	currSet := make(map[scanner.Port]struct{}, len(curr))
	for _, p := range curr {
		currSet[p] = struct{}{}
	}
	for _, p := range curr {
		if _, ok := prevSet[p]; !ok {
			added = append(added, p)
		}
	}
	for _, p := range prev {
		if _, ok := currSet[p]; !ok {
			removed = append(removed, p)
		}
	}
	return
}
