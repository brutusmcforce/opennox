package eval

import "reflect"

// Scripts import these packages by their public paths, which are the keys in Symbols.
// The packages now live in OpenNox (github.com/brutusmcforce/opennox/...), so the extracted
// files must be edited to keep these keys and to import the in-repo packages.
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4/abil
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4/audio
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4/class
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4/damage
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4/effect
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4/enchant
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4/spell
//go:generate yaegi extract github.com/opennox/noxscript/ns/v4/subclass

//go:generate goimports -w .

var Symbols = make(map[string]map[string]reflect.Value)
