package layout

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

var testBuf [64]byte

func TestCheckAs(t *testing.T) {
	old := Checking
	Checking = true
	defer func() { Checking = old }()
	type inner struct{ A, B int32 }
	type outer struct {
		X  int32
		In inner
	}
	p := unsafe.Pointer(&testBuf[0]) // not on the stack: tracked addresses must not move
	Track(p, unsafe.Sizeof(outer{}), reflect.TypeFor[outer](), "")
	defer Untrack(p)

	reported.seen = nil
	var got []string
	reportHook = func(s string) { got = append(got, s) }
	defer func() { reportHook = nil }()

	As[outer](p)
	As[inner](unsafe.Add(p, unsafe.Offsetof(outer{}.In)))
	As[int32](unsafe.Add(p, unsafe.Offsetof(outer{}.In)+4))
	if len(got) != 0 {
		t.Fatalf("unexpected reports: %q", got)
	}
	As[inner](p)
	if len(got) != 1 || !strings.Contains(got[0], "WRONG TYPE") {
		t.Fatalf("expected a wrong type report, got %q", got)
	}
}

func TestCheckAsUnion(t *testing.T) {
	old := Checking
	Checking = true
	defer func() { Checking = old }()
	type union struct{ _ [2]uint32 }
	type varA struct{ A, B uint32 }
	type varB struct{ C uint16 }
	type other struct{ D uint32 }
	type outer struct {
		X uint32
		U union
	}
	Union[union](varA{}, varB{})
	p := unsafe.Pointer(&testBuf[0])
	Track(p, unsafe.Sizeof(outer{}), reflect.TypeFor[outer](), "")
	defer Untrack(p)

	reported.seen = nil
	var got []string
	reportHook = func(s string) { got = append(got, s) }
	defer func() { reportHook = nil }()

	u := unsafe.Add(p, unsafe.Offsetof(outer{}.U))
	As[varA](u)
	As[varB](u)
	if len(got) != 0 {
		t.Fatalf("unexpected reports: %q", got)
	}
	As[other](u)
	As[varA](unsafe.Add(u, 4))
	if len(got) != 2 {
		t.Fatalf("expected two wrong type reports, got %q", got)
	}
}

type rawTestS struct {
	A uint32
	B uint16
	P *int
	C uint32
}

var (
	rawTestStruct rawTestS // not on the stack: tracked addresses must not move
	rawTestBuf    [32]byte
)

func TestCheckRaw(t *testing.T) {
	old := Checking
	Checking = true
	defer func() { Checking = old }()
	var got []string
	reportHook = func(s string) { got = append(got, s) }
	defer func() { reportHook = nil }()
	s := &rawTestStruct
	TrackAs[rawTestS](unsafe.Pointer(s))
	defer Untrack(unsafe.Pointer(s))
	raw := rawTestBuf[:]
	Track(unsafe.Pointer(&raw[0]), 32, nil, "raw")
	defer Untrack(unsafe.Pointer(&raw[0]))

	CheckRaw(unsafe.Pointer(s), 0, 4, 0)        // A: same place
	CheckRaw(unsafe.Pointer(s), 4, 2, 0)        // B: same place
	CheckRaw(unsafe.Pointer(&raw[0]), 12, 4, 0) // untyped block
	CheckRaw(unsafe.Pointer(&s.B), 0, 2, 0)     // interior base
	if len(got) != 0 {
		t.Fatalf("unexpected reports: %q", got)
	}
	CheckRaw(unsafe.Pointer(s), 12, 4, 0) // C: 12 on 386, 16 on amd64
	if runtime.GOARCH == "amd64" && len(got) != 1 {
		t.Fatalf("got %q, want one report for C", got)
	}
}
