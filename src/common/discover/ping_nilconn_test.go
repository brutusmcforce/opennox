package discover

import (
	"errors"
	"log/slog"
	"net/netip"
	"testing"
	"time"

	"github.com/brutusmcforce/opennox/lobby"
)

// The broadcast socket cannot always be opened (a firewall, a port in use, or a Windows API that the host
// does not implement). The pinger must then report it instead of using the nil connection: it used to start
// its read loop anyway and crash the client with a nil pointer dereference.
func TestPingerWithoutSocket(t *testing.T) {
	p := NewPinger(slog.Default(), nil)

	err := p.SendPing(make(chan<- *lobby.Game, 1), netip.MustParseAddrPort("127.0.0.1:18700"))
	if !errors.Is(err, ErrNoPingSocket) {
		t.Fatalf("SendPing: got %v, want %v", err, ErrNoPingSocket)
	}

	done := make(chan error, 1)
	go func() { done <- p.Close() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Close: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Close blocked")
	}
	if err := p.Close(); err != nil { // a second Close must still be a no-op
		t.Fatalf("Close twice: %v", err)
	}
}
