package legacy

import "unsafe"

import (
	"bytes"
	"context"
	"strings"

	"github.com/brutusmcforce/opennox/libs/ifs"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
)

// commandRuleHeader is the exact narrow C strcmp lookup used by the command
// rule reader. It does not mutate the settings parser's ruleLoaderContext.
func commandRuleHeader(s string) uint32 {
	s = ruleCString(s)
	for i := 0; i < 7; i++ {
		p := (*int8)(*memmap.PtrPtr(0x587000, uintptr(312208+8*i)))
		if p != nil && GoStringP(unsafe.Pointer(p)) == s {
			return *memmap.PtrUint32(0x587000, uintptr(312212+8*i))
		}
	}
	return 0
}

// commandRulesFile reads a command rule file using the existing text-file
// bridge semantics: one complete physical line is read, CRLF is normalized by
// binfile, then fgets' 255-byte destination keeps at most 254 bytes.
func commandRulesFile(path string) int {
	f, err := ifs.Open(ruleCString(path))
	if err != nil {
		return 0
	}
	bf := binfile.NewTextFile(f)
	defer bf.Close()
	section := uint32(0x17f0)
	for {
		line, readErr := bf.ReadString()
		if len(line) != 0 {
			if len(line) > 254 {
				line = line[:254]
			}
			if i := bytes.IndexByte(line, 0); i >= 0 {
				line = line[:i]
			}
			if i := bytes.IndexByte(line, '\n'); i >= 0 {
				line = line[:i]
			}
			if len(line) != 0 {
				if mode := commandRuleHeader(string(line)); mode != 0 {
					section = mode
				} else if noxflags.HasGame(noxflags.GameFlag(section)) {
					// C widens each raw byte before invoking the Go-backed command parser.
					wide := make([]rune, len(line))
					for i, b := range line {
						wide[i] = rune(b)
					}
					_ = ExecConsoleCmd(context.Background(), string(wide))
				}
			}
		}
		if readErr != nil {
			// EOF can accompany a final physical line, which was processed above.
			// The C bridge can spin on a non-EOF read failure; terminate safely.
			break
		}
	}
	return 1
}

// commandRulesPath preserves the nullable C entry point. It applies the C
// input-buffer limit before extension removal and recognizes only backslash as
// a directory separator while constructing the first user.rul probe.
func commandRulesPath(path *string) int {
	if path == nil {
		return 0
	}
	p := ruleCString(*path)
	if len(p) > 255 {
		p = p[:255]
	}
	if len(p) < 4 {
		return 0
	}
	stem := p[:len(p)-4]
	dir := stem
	if i := strings.LastIndexByte(stem, '\\'); i >= 0 {
		dir = stem[:i+1]
	}
	if commandRulesFile(dir+"user.rul") == 0 {
		_ = commandRulesFile(stem + ".rul")
	}
	return 1
}

func commandRulesMap(name string) int {
	name = ruleCString(name)
	prefix := "maps\\" + name
	if len(prefix) > 255 {
		prefix = prefix[:255]
	}
	path := prefix[:len(prefix)-4] + "\\" + name
	return commandRulesPath(&path)
}
