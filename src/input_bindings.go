package opennox

import (
	"strings"

	"github.com/opennox/libs/client/keybind"
)

// OpenNox extensions to the key binding events. Unlike the ones in keybind package, they are handled only by OpenNox.
const (
	eventZoomIn = keybind.Event(0x100 + iota)
	eventZoomOut
)

var extBindEvents = []keybind.BindEvent{
	{Name: "ZoomIn", Event: eventZoomIn, Title: "Zoom Camera In"},
	{Name: "ZoomOut", Event: eventZoomOut, Title: "Zoom Camera Out"},
}

// defaultExtBindings are OpenNox changes to the default key bindings from default.cfg.
var defaultExtBindings = []struct {
	key keybind.Key
	ev  keybind.Event
}{
	{keybind.KeyMouseWheelUp, eventZoomIn},
	{keybind.KeyMouseWheelDown, eventZoomOut},
}

// bindEvents is a list of all key binding events, including OpenNox extensions.
type bindEvents struct {
	base *keybind.Binding
	list []keybind.BindEvent
}

func newBindEvents(base *keybind.Binding) *bindEvents {
	b := &bindEvents{base: base}
	b.list = append(b.list, base.Events()...)
	b.list = append(b.list, extBindEvents...)
	return b
}

// Events returns all events, except the empty one.
func (b *bindEvents) Events() []keybind.BindEvent {
	return b.list
}

func (b *bindEvents) EventByCode(ev keybind.Event) *keybind.BindEvent {
	for i := range extBindEvents {
		if e := &extBindEvents[i]; e.Event == ev {
			return e
		}
	}
	return b.base.EventByCode(ev)
}

func (b *bindEvents) EventByName(name string) *keybind.BindEvent {
	for i := range extBindEvents {
		if e := &extBindEvents[i]; e.Name == name {
			return e
		}
	}
	return b.base.EventByName(name)
}

func (b *bindEvents) EventByTitle(title string) *keybind.BindEvent {
	for i := range extBindEvents {
		if e := &extBindEvents[i]; strings.EqualFold(e.Title, title) {
			return e
		}
	}
	return b.base.EventByTitle(title)
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
