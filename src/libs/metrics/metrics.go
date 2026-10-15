// Package metrics collects the counters and gauges of the game and serves them in the Prometheus text
// exposition format.
//
// The game only ever needs counters and gauges with a few labels, which is little enough to keep here: the
// Prometheus client library costs 3.3 MB in the binary and brings in eight more modules (protobuf, procfs,
// compress...) for it. The lobby server does more and still uses the real client library.
//
// What is lost by not using it are the go_* and process_* collectors it registers by default; add them here
// (from runtime/metrics) if a server ever needs them.
package metrics

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// Opts are the options of a metric. The name must be a valid Prometheus metric name.
type Opts struct {
	Name string
	Help string
}

type kind string

const (
	kindCounter = kind("counter")
	kindGauge   = kind("gauge")
)

// value is a float64 that can be updated from any goroutine.
type value struct {
	bits atomic.Uint64
}

func (v *value) add(d float64) {
	for {
		old := v.bits.Load()
		if v.bits.CompareAndSwap(old, math.Float64bits(math.Float64frombits(old)+d)) {
			return
		}
	}
}

func (v *value) set(f float64) { v.bits.Store(math.Float64bits(f)) }
func (v *value) get() float64  { return math.Float64frombits(v.bits.Load()) }

// Metric is a single counter or gauge value. Counters are only ever increased, gauges are also set.
type Metric struct {
	value
}

func (m *Metric) Inc()           { m.add(1) }
func (m *Metric) Dec()           { m.add(-1) }
func (m *Metric) Add(d float64)  { m.add(d) }
func (m *Metric) Set(f float64)  { m.set(f) }
func (m *Metric) Value() float64 { return m.get() }

// Vec is a counter or gauge split by label values.
type Vec struct {
	labels []string

	mu     sync.Mutex
	byKey  map[string]*Metric
	values map[string][]string
}

// WithLabelValues returns the metric for one set of label values, creating it on first use.
// It panics if the number of values does not match the labels, like the Prometheus client does.
func (v *Vec) WithLabelValues(vals ...string) *Metric {
	if len(vals) != len(v.labels) {
		panic(fmt.Errorf("metrics: got %d label values, expected %d", len(vals), len(v.labels)))
	}
	key := strings.Join(vals, "\x00")
	v.mu.Lock()
	defer v.mu.Unlock()
	m := v.byKey[key]
	if m == nil {
		m = new(Metric)
		v.byKey[key] = m
		v.values[key] = append([]string(nil), vals...)
	}
	return m
}

// DeleteLabelValues stops reporting the metric for one set of label values. It reports whether it was
// there. The lobby uses it to drop the series of a game that went offline, which would otherwise be
// reported forever.
func (v *Vec) DeleteLabelValues(vals ...string) bool {
	if len(vals) != len(v.labels) {
		return false
	}
	key := strings.Join(vals, "\x00")
	v.mu.Lock()
	defer v.mu.Unlock()
	if _, ok := v.byKey[key]; !ok {
		return false
	}
	delete(v.byKey, key)
	delete(v.values, key)
	return true
}

// collector is one registered metric family.
type collector struct {
	name string
	help string
	kind kind
	m    *Metric // set for a plain metric
	vec  *Vec    // set for a labelled one
}

var reg struct {
	mu  sync.Mutex
	all []*collector
}

func register(c *collector) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.all = append(reg.all, c)
}

func newMetric(k kind, o Opts) *Metric {
	m := new(Metric)
	register(&collector{name: o.Name, help: o.Help, kind: k, m: m})
	return m
}

func newVec(k kind, o Opts, labels []string) *Vec {
	v := &Vec{
		labels: append([]string(nil), labels...),
		byKey:  make(map[string]*Metric),
		values: make(map[string][]string),
	}
	register(&collector{name: o.Name, help: o.Help, kind: k, vec: v})
	return v
}

// NewCounter creates a counter and registers it.
func NewCounter(o Opts) *Metric { return newMetric(kindCounter, o) }

// NewGauge creates a gauge and registers it.
func NewGauge(o Opts) *Metric { return newMetric(kindGauge, o) }

// NewCounterVec creates a counter split by the given labels and registers it.
func NewCounterVec(o Opts, labels []string) *Vec { return newVec(kindCounter, o, labels) }

// NewGaugeVec creates a gauge split by the given labels and registers it.
func NewGaugeVec(o Opts, labels []string) *Vec { return newVec(kindGauge, o, labels) }

// Escaping of the text format: label values are quoted, the help text is not.
var (
	escapeLabel = strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`).Replace
	escapeHelp  = strings.NewReplacer(`\`, `\\`, "\n", `\n`).Replace
)

// formatValue writes a float the way Prometheus expects it.
func formatValue(f float64) string {
	switch {
	case math.IsInf(f, 1):
		return "+Inf"
	case math.IsInf(f, -1):
		return "-Inf"
	case math.IsNaN(f):
		return "NaN"
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// Write writes all metrics in the Prometheus text exposition format.
func Write(w io.Writer) error {
	reg.mu.Lock()
	all := append([]*collector(nil), reg.all...)
	reg.mu.Unlock()

	var buf strings.Builder
	for _, c := range all {
		if c.help != "" {
			fmt.Fprintf(&buf, "# HELP %s %s\n", c.name, escapeHelp(c.help))
		}
		fmt.Fprintf(&buf, "# TYPE %s %s\n", c.name, c.kind)
		if c.m != nil {
			fmt.Fprintf(&buf, "%s %s\n", c.name, formatValue(c.m.Value()))
			continue
		}
		v := c.vec
		v.mu.Lock()
		keys := make([]string, 0, len(v.byKey))
		for key := range v.byKey {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			vals := v.values[key]
			var pairs strings.Builder
			for i, name := range v.labels {
				if i != 0 {
					pairs.WriteByte(',')
				}
				// not %q: that would also escape the non-ASCII characters of e.g. a player name,
				// which the text format keeps as UTF-8
				fmt.Fprintf(&pairs, `%s="%s"`, name, escapeLabel(vals[i]))
			}
			fmt.Fprintf(&buf, "%s{%s} %s\n", c.name, pairs.String(), formatValue(v.byKey[key].Value()))
		}
		v.mu.Unlock()
	}
	_, err := io.WriteString(w, buf.String())
	return err
}

// Handler serves the metrics, like the Prometheus client's promhttp.Handler.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		_ = Write(w)
	})
}
