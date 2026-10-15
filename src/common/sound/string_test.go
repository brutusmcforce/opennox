package sound

import "testing"

func TestIDStringOutOfRange(t *testing.T) {
	for _, id := range []ID{-1, ID(len(soundNames)), 1 << 20} {
		if got, want := id.String(), "sound.ID("; len(got) < len(want) || got[:len(want)] != want {
			t.Errorf("%d: got %q", int(id), got)
		}
	}
	for _, id := range []ID{-1, ID(len(soundNames))} {
		if got, want := id.GoString(), "sound.ID("; len(got) < len(want) || got[:len(want)] != want {
			t.Errorf("%d: got %q", int(id), got)
		}
	}
	if got := SoundWalkOnStone.String(); got != "WalkOnStone" {
		t.Errorf("got %q", got)
	}
}
