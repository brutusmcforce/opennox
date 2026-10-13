package xfer

// StaticRegistry stores static XFER decoder mappings for Nox objects. See DefaultRegistry.
type StaticRegistry struct {
	ByType map[string]Type
	ByID   map[int]Type
	// Classes are the class and subclass bits of the object types whose XFER data depends on them.
	Classes map[string]ObjectClass
}

func (r *StaticRegistry) XferByObjectType(typ string) Type {
	return r.ByType[typ]
}

func (r *StaticRegistry) XferByObjectTypeID(id int) Type {
	return r.ByID[id]
}

// ObjectClassByType implements ClassRegistry.
func (r *StaticRegistry) ObjectClassByType(typ string) (ObjectClass, bool) {
	c, ok := r.Classes[typ]
	return c, ok
}

// ObjectClass holds the class and subclass bits of an object type (as in thing.bin).
type ObjectClass struct {
	Class    uint32
	SubClass uint32
}

// ClassRegistry is an optional ObjectRegistry interface for the class bits of object types. Some XFER data depends
// on them: wand charges, and a monster's shop, colors and voice.
type ClassRegistry interface {
	ObjectClassByType(typ string) (ObjectClass, bool)
}

// TypeNameRegistry is an optional ObjectRegistry interface that finds object type names by ID.
type TypeNameRegistry interface {
	ObjectTypeByID(id int) string
}

func objectClass(reg ObjectRegistry, typ string) ObjectClass {
	if cr, ok := reg.(ClassRegistry); ok && typ != "" {
		c, _ := cr.ObjectClassByType(typ)
		return c
	}
	return ObjectClass{}
}

// TOCRegistry finds object types by the indexes of a map's object TOC (as used by map object data), with XFER types
// and classes from Base.
type TOCRegistry struct {
	Base ObjectRegistry
	TOC  map[int]string
}

func (r *TOCRegistry) XferByObjectType(typ string) Type {
	return r.Base.XferByObjectType(typ)
}

func (r *TOCRegistry) XferByObjectTypeID(id int) Type {
	typ, ok := r.TOC[id]
	if !ok {
		return ""
	}
	return r.Base.XferByObjectType(typ)
}

func (r *TOCRegistry) ObjectTypeByID(id int) string {
	return r.TOC[id]
}

func (r *TOCRegistry) ObjectClassByType(typ string) (ObjectClass, bool) {
	if cr, ok := r.Base.(ClassRegistry); ok {
		return cr.ObjectClassByType(typ)
	}
	return ObjectClass{}, false
}
