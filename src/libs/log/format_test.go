package log

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

// writesRecorder keeps every Write call separately.
type writesRecorder struct {
	mu     sync.Mutex
	writes []string
}

func (w *writesRecorder) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writes = append(w.writes, string(p))
	return len(p), nil
}

// Records logged from several goroutines must each be written whole, with a single call.
func TestTextHandlerConcurrent(t *testing.T) {
	var w writesRecorder
	h := NewTextHandler(&w)
	const (
		workers = 8
		each    = 200
	)
	var wg sync.WaitGroup
	for i := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range each {
				r := slog.NewRecord(time.Now(), slog.LevelInfo, "message", 0)
				r.AddAttrs(slog.Int("worker", i), slog.Int("n", j))
				if err := h.Handle(context.Background(), r); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Wait()
	if len(w.writes) != workers*each {
		t.Fatalf("got %d writes, want %d", len(w.writes), workers*each)
	}
	for _, s := range w.writes {
		if strings.Count(s, "\n") != 1 || !strings.HasSuffix(s, "\n") || !strings.Contains(s, "INFO [main] message worker=") {
			t.Fatalf("not a whole record: %q", s)
		}
	}
}

// The log file used to be written through a bufio.Writer, which isn't safe for concurrent use.
func TestTextHandlerConcurrentBuffered(t *testing.T) {
	var out bytes.Buffer
	bw := bufio.NewWriter(&out)
	h := NewTextHandler(bw)
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 200 {
				r := slog.NewRecord(time.Now(), slog.LevelInfo, "message", 0)
				r.AddAttrs(slog.Int("worker", i), slog.Int("n", j))
				_ = h.Handle(context.Background(), r)
			}
		}()
	}
	wg.Wait()
	if err := bw.Flush(); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if len(lines) != 8*200 {
		t.Fatalf("got %d lines", len(lines))
	}
	for _, l := range lines {
		if !strings.Contains(l, "INFO [main] message worker=") {
			t.Fatalf("broken line: %q", l)
		}
	}
}

func TestPrintRecord(t *testing.T) {
	r := slog.NewRecord(time.Date(2026, 10, 15, 8, 30, 0, 0, time.UTC), slog.LevelWarn, "hello\n", 0)
	r.AddAttrs(slog.String(systemAttr, "net"), slog.Int("n", 1))
	var buf bytes.Buffer
	if err := PrintRecord(&buf, PrintDefault, r); err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), "2026-10-15 08:30:00 WARN [net] hello n=1\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
