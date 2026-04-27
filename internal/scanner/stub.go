package scanner

// Scannable is the interface satisfied by Scanner and test stubs.
type Scannable interface {
	Scan() ([]Port, error)
}

// StubBackend is implemented by types that provide canned scan results.
type StubBackend interface {
	Scan() ([]Port, error)
}

// stubScanner wraps any StubBackend so it satisfies the same interface
// as *Scanner without exposing internal fields.
type stubScanner struct {
	backend StubBackend
}

// WrapStub wraps a StubBackend as a *Scanner-compatible Scannable.
// Used exclusively in tests.
func WrapStub(b StubBackend) *Scanner {
	// We embed the backend call inside a real Scanner whose dial func
	// is never used; instead we override the Scan method via composition.
	// Since Scanner.Scan is the only public method we care about, we
	// replace the internal dialer with one driven by the stub.
	s := &Scanner{}
	s.dialFn = func(network, address string) error {
		// unused when backend is set
		return nil
	}
	s.stub = b
	return s
}
