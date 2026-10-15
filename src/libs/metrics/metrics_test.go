package metrics

import (
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestWriteTextFormat(t *testing.T) {
	cnt := NewCounter(Opts{Name: "test_total", Help: "A counter"})
	g := NewGauge(Opts{Name: "test_gauge", Help: "A gauge"})
	vec := NewCounterVec(Opts{Name: "test_by_map", Help: "By map"}, []string{"map"})

	cnt.Inc()
	cnt.Inc()
	g.Set(2.5)
	vec.WithLabelValues("estate").Inc()
	vec.WithLabelValues(`we"ird\one`).Inc()

	var buf strings.Builder
	if err := Write(&buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		"# HELP test_total A counter\n# TYPE test_total counter\ntest_total 2\n",
		"# TYPE test_gauge gauge\ntest_gauge 2.5\n",
		`test_by_map{map="estate"} 1` + "\n",
		`test_by_map{map="we\"ird\\one"} 1` + "\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output does not contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestHandler(t *testing.T) {
	NewGauge(Opts{Name: "handler_gauge", Help: "H"}).Set(7)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest("GET", "/metrics", nil))
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("content type %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "handler_gauge 7\n") {
		t.Errorf("body: %s", rec.Body.String())
	}
}

// The metrics are updated from the server tick and read by a scrape at the same time.
func TestConcurrentUpdates(t *testing.T) {
	cnt := NewCounter(Opts{Name: "race_total", Help: "R"})
	vec := NewGaugeVec(Opts{Name: "race_by_id", Help: "R"}, []string{"id"})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				cnt.Inc()
				vec.WithLabelValues("a").Set(float64(j))
			}
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100; j++ {
			var buf strings.Builder
			_ = Write(&buf)
		}
	}()
	wg.Wait()
	if got := cnt.Value(); got != 8000 {
		t.Errorf("counter = %v, want 8000", got)
	}
}

func TestWrongLabelCountPanics(t *testing.T) {
	vec := NewCounterVec(Opts{Name: "panic_by_two", Help: "P"}, []string{"a", "b"})
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic")
		}
	}()
	vec.WithLabelValues("only-one")
}

func TestDeleteLabelValues(t *testing.T) {
	vec := NewGaugeVec(Opts{Name: "del_by_game", Help: "D"}, []string{"addr"})
	vec.WithLabelValues("a").Set(1)
	vec.WithLabelValues("b").Set(2)

	if !vec.DeleteLabelValues("a") {
		t.Error("DeleteLabelValues(a) = false, want true")
	}
	if vec.DeleteLabelValues("a") {
		t.Error("deleting twice reported true")
	}
	if vec.DeleteLabelValues("x", "y") { // wrong label count
		t.Error("wrong label count reported true")
	}

	var buf strings.Builder
	if err := Write(&buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, `del_by_game{addr="a"}`) {
		t.Errorf("deleted series still reported:\n%s", out)
	}
	if !strings.Contains(out, `del_by_game{addr="b"} 2`) {
		t.Errorf("remaining series missing:\n%s", out)
	}
}
