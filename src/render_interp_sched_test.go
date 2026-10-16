package opennox

import (
	"testing"
	"time"
)

// The frame limiter must keep the frame rate at the limit, with the game tick frames counted in.
func TestFrameScheduleRate(t *testing.T) {
	const tick = time.Second / 30
	for _, fps := range []int{60, 75, 90, 100, 120, 144, 165, 180, 200, 240, 360} {
		period := time.Second / time.Duration(fps)
		var next, now time.Duration
		frames := 0
		for n := 1; n <= 30*10; n++ { // 10 seconds
			// a tick frame is drawn on each game tick
			now = time.Duration(n) * tick
			next = nextFrameTime(next, now, period)
			frames++
			// additional frames until the next tick, each drawn when it's due (see waitAndRender)
			deadline := time.Duration(n+1) * tick
			for next < deadline {
				now = max(next, now)
				next = nextFrameTime(next, now, period)
				frames++
			}
		}
		got := float64(frames) / 10
		if got < float64(fps)*0.99 || got > float64(fps)*1.01 {
			t.Errorf("limit %d fps: %.1f fps", fps, got)
		}
	}
}
