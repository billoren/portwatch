package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func makePort(proto string, num uint16) scanner.Port {
	return scanner.Port{Proto: proto, Number: num}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	ports := []scanner.Port{
		makePort("tcp", 80),
		makePort("tcp", 443),
	}

	if err := state.SaveSnapshot(path, ports); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	snap, err := state.LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if len(snap.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadSnapshotMissingFile(t *testing.T) {
	snap, err := state.LoadSnapshot("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if snap.Ports != nil {
		t.Error("expected nil ports for missing snapshot")
	}
}

func TestLoadSnapshotCorrupt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := state.LoadSnapshot(path)
	if err == nil {
		t.Fatal("expected error for corrupt snapshot")
	}
}

func TestDiff(t *testing.T) {
	prev := []scanner.Port{makePort("tcp", 80), makePort("tcp", 22)}
	curr := []scanner.Port{makePort("tcp", 80), makePort("tcp", 8080)}

	added, removed := state.Diff(prev, curr)

	if len(added) != 1 || added[0] != makePort("tcp", 8080) {
		t.Errorf("unexpected added: %v", added)
	}
	if len(removed) != 1 || removed[0] != makePort("tcp", 22) {
		t.Errorf("unexpected removed: %v", removed)
	}
}

func TestDiffNoChanges(t *testing.T) {
	ports := []scanner.Port{makePort("tcp", 443)}
	added, removed := state.Diff(ports, ports)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}
