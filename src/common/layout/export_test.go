package layout

import "reflect"

var (
	Size386    = size386
	Build      = build
	Model386   = model386
	ModelAmd64 = modelAmd64
)

func (m model) Size(t reflect.Type) uintptr               { return m.size(t) }
func (m model) OffsetOf(t reflect.Type, p string) uintptr { return m.offsetOf(t, p) }
func (t *Table) Offsets() []uint32                        { return t.off }

var FieldName = fieldName
