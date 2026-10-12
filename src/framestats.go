package opennox

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/brutusmcforce/opennox/libs/log"
)

// Frame statistics, enabled with NOX_FRAMESTATS: "1" logs a summary of frame times when the game exits
// (or the E2E script completes), any other value is a CSV file that also gets the time of every frame.
// Times are measured in real time between presented frames, so they include everything the game does in a frame.

var (
	frameStatsLog = log.New("framestats")
	frameStats    = newFrameStatsRecorder(os.Getenv("NOX_FRAMESTATS"))
)

// frameKind tells what a presented frame showed, so that menus and loading screens can be told apart from the game.
type frameKind int

const (
	frameMenu frameKind = iota
	frameGame
	frameMovie
	frameHarness // frames with test harness work (E2E screenshots), not counted as game frames
	frameKinds
)

func (k frameKind) String() string {
	switch k {
	case frameGame:
		return "game"
	case frameMovie:
		return "movie"
	case frameHarness:
		return "e2e"
	default:
		return "menu"
	}
}

// slowFrame is the frame time above which game frames are logged.
const slowFrame = 12 * time.Millisecond

type frameStatsRecorder struct {
	on           bool
	csv          *os.File
	w            *bufio.Writer
	last         time.Time
	times        [frameKinds][]time.Duration
	harnessFrame bool
}

// harness marks the current frame as one with test harness work.
func (s *frameStatsRecorder) harness() {
	s.harnessFrame = true
}

func newFrameStatsRecorder(conf string) *frameStatsRecorder {
	s := &frameStatsRecorder{on: conf != ""}
	if conf != "" && conf != "1" {
		f, err := os.Create(conf)
		if err != nil {
			frameStatsLog.Printf("cannot create %q: %v", conf, err)
		} else {
			s.csv = f
			s.w = bufio.NewWriter(f)
			fmt.Fprintln(s.w, "ms,kind")
		}
	}
	return s
}

// presented records a presented frame.
func (s *frameStatsRecorder) presented(kind frameKind) {
	if !s.on {
		return
	}
	now := time.Now()
	if s.harnessFrame {
		s.harnessFrame = false
		kind = frameHarness
	}
	if !s.last.IsZero() {
		dt := now.Sub(s.last)
		s.times[kind] = append(s.times[kind], dt)
		if kind == frameGame && dt > slowFrame {
			// logged in order with the rest of the log, which shows what the game was doing
			frameStatsLog.Printf("slow frame %d: %.1f ms", len(s.times[kind]), float64(dt)/float64(time.Millisecond))
		}
		if s.w != nil {
			fmt.Fprintf(s.w, "%.3f,%s\n", float64(dt)/float64(time.Millisecond), kind)
		}
	}
	s.last = now
}

// finish stops recording and logs the summary.
func (s *frameStatsRecorder) finish() {
	if !s.on {
		return
	}
	s.on = false
	for k := frameMenu; k < frameKinds; k++ {
		if st, ok := summarizeFrames(s.times[k]); ok {
			frameStatsLog.Printf("%s: %s", k, st)
		}
	}
	if s.w != nil {
		_ = s.w.Flush()
		_ = s.csv.Close()
	}
}

// frameSummary describes frame times in frames per second.
type frameSummary struct {
	Frames int
	AvgFPS float64 // frames divided by the total time
	Low1   float64 // 1% low: the frame rate of the slowest 1% of frames (99th percentile of frame times)
	Worst  time.Duration
}

func (s frameSummary) String() string {
	return fmt.Sprintf("%d frames, avg %.0f fps, 1%% low %.0f fps, worst frame %.1f ms (%.0f fps)",
		s.Frames, s.AvgFPS, s.Low1, float64(s.Worst)/float64(time.Millisecond), float64(time.Second)/float64(s.Worst))
}

func summarizeFrames(times []time.Duration) (frameSummary, bool) {
	if len(times) == 0 {
		return frameSummary{}, false
	}
	var total time.Duration
	for _, dt := range times {
		total += dt
	}
	sorted := slices.Clone(times)
	slices.Sort(sorted)
	p99 := sorted[min(len(sorted)-1, len(sorted)*99/100)]
	return frameSummary{
		Frames: len(times),
		AvgFPS: float64(len(times)) * float64(time.Second) / float64(total),
		Low1:   float64(time.Second) / float64(max(p99, 1)),
		Worst:  sorted[len(sorted)-1],
	}, true
}
