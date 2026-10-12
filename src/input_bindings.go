package opennox

import (
	"github.com/brutusmcforce/opennox/libs/client/keybind"
)

// defaultExtBindings are OpenNox changes to the default key bindings from default.cfg.
var defaultExtBindings = []struct {
	key keybind.Key
	ev  keybind.Event
}{
	{keybind.KeyMouseWheelUp, keybind.EventZoomIn},
	{keybind.KeyMouseWheelDown, keybind.EventZoomOut},
}

// applyDefaultExtBindings applies OpenNox changes to the default key bindings.
func (c *CtrlEventHandler) applyDefaultExtBindings() {
	for _, b := range defaultExtBindings {
		c.setKeyBinding(b.key, b.ev)
	}
}

// setKeyBinding binds a single key to an event, replacing all other bindings of this key.
func (c *CtrlEventHandler) setKeyBinding(key keybind.Key, ev keybind.Event) {
	for it := c.bindings; it != nil; {
		next := it.next
		if len(it.keys) == 1 && it.keys[0] == key {
			c.removeBinding(it)
		}
		it = next
	}
	c.addBinding(&CtrlEventBinding{keys: []keybind.Key{key}, events: []keybind.Event{ev}})
}

func (c *CtrlEventHandler) removeBinding(ce *CtrlEventBinding) {
	if ce.prev != nil {
		ce.prev.next = ce.next
	} else if c.bindings == ce {
		c.bindings = ce.next
	}
	if ce.next != nil {
		ce.next.prev = ce.prev
	}
	ce.prev, ce.next = nil, nil
}
