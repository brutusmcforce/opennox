package keybind

import "testing"

func TestBindEventsUnique(t *testing.T) {
	names := make(map[string]bool)
	codes := make(map[Event]bool)
	for _, e := range defBindEvents[1:] {
		if names[e.Name] {
			t.Errorf("duplicate event name: %q", e.Name)
		}
		names[e.Name] = true
		if codes[e.Event] {
			t.Errorf("duplicate event code: %v (%q)", e.Event, e.Name)
		}
		codes[e.Event] = true
	}
	for _, name := range []string{"ZoomIn", "ZoomOut"} {
		if !names[name] {
			t.Errorf("missing event: %q", name)
		}
	}
}
