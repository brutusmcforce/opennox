package legacy

import (
	"bytes"
	"unicode/utf16"
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

type wchar2_t = uint16

func StrFree[T comparable](s *T) {
	legacyFree(unsafe.Pointer(s))
}

func StrLenBytes(s []byte) int {
	i := bytes.IndexByte(s, 0)
	if i < 0 {
		return len(s)
	}
	return i
}

func StrCopy(dst *int8, max int, src string) int {
	d := unsafe.Slice((*byte)(unsafe.Pointer(dst)), max)
	return StrCopyBytes(d, src)
}

func StrNCopy(dst *int8, max int, src string) int {
	d := unsafe.Slice((*byte)(unsafe.Pointer(dst)), max)
	return StrNCopyBytes(d, src)
}

func StrCopyP(dst unsafe.Pointer, max int, src string) int {
	return StrCopy((*int8)(dst), max, src)
}

func StrCopyBytes(dst []byte, src string) int {
	n := copy(dst[:len(dst)-1], src)
	dst[n] = 0
	return n
}

func StrNCopyBytes(dst []byte, src string) int {
	n := copy(dst, src)
	if n < len(dst) {
		dst[n] = 0
	}
	return n
}

func WStrLenBytes(s []byte) int {
	if s == nil {
		return 0
	}
	n := 0
	for len(s) >= 2 && s[0] != 0 && s[1] != 0 {
		s = s[2:]
		n++
	}
	return n
}

func WStrCopy(p *wchar2_t, max int, src string) int {
	dst := unsafe.Slice((*uint16)(unsafe.Pointer(p)), max)
	val := utf16.Encode([]rune(src))
	n := copy(dst, val)
	if n < len(dst) {
		dst[n] = 0
	}
	return n
}

func WStrCopyBytes(p []byte, src string) int {
	dst := unsafe.Slice((*uint16)(unsafe.Pointer(&p[0])), len(p)/2)
	val := utf16.Encode([]rune(src))
	n := copy(dst, val)
	if n < len(dst) {
		dst[n] = 0
	}
	return n
}

func WStrCopySlice(dst []uint16, src string) int {
	val := utf16.Encode([]rune(src))
	n := copy(dst, val)
	if n < len(dst) {
		dst[n] = 0
	}
	return n
}

func GoString(s *int8) string {
	if s == nil {
		return ""
	}
	n := 0
	for p := (*byte)(unsafe.Pointer(s)); *p != 0; p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
		n++
	}
	return string(unsafe.Slice((*byte)(unsafe.Pointer(s)), n))
}

func GoStringP(s unsafe.Pointer) string {
	return GoString((*int8)(s))
}

func GoStringN(s *int8, n int) string {
	return GoStringNP(unsafe.Pointer(s), n)
}

func GoStringNP(s unsafe.Pointer, n int) string {
	b := unsafe.Slice((*byte)(s), n)
	i := bytes.IndexByte(b, 0)
	if i < 0 {
		return string(b)
	}
	return string(b[:i])
}

func GoStringS(s []byte) string {
	return string(s[:StrLenBytes(s)])
}

func CString(s string) *int8 {
	if len(s)+1 <= 0 {
		panic("string too large")
	}
	p := legacyMalloc(uintptr(len(s) + 1))
	buf := unsafe.Slice((*byte)(p), len(s)+1)
	copy(buf, s)
	buf[len(s)] = 0
	return (*int8)(p)
}

func GoWStringP(s unsafe.Pointer) string {
	return GoWString((*wchar2_t)(s))
}

func GoWString(s *wchar2_t) string {
	n := alloc.StrLen(s)
	if n == 0 {
		return ""
	}
	arr := unsafe.Slice((*uint16)(unsafe.Pointer(s)), n)
	return GoWStringSlice(arr)
}

func GoWStringN(s *wchar2_t, max int) string {
	n := alloc.StrLenN(s, max)
	if n == 0 {
		return ""
	}
	arr := unsafe.Slice((*uint16)(unsafe.Pointer(s)), n)
	return GoWStringSlice(arr)
}

func GoWStringSlice(arr []uint16) string {
	return string(utf16.Decode(arr))
}

func GoWStringBytes(arr []byte) string {
	n := WStrLenBytes(arr)
	if n == 0 {
		return ""
	}
	warr := unsafe.Slice((*uint16)(unsafe.Pointer(&arr[0])), n)
	return GoWStringSlice(warr)
}

func CWString(s string) (*wchar2_t, func()) {
	buf := utf16.Encode([]rune(s))
	ptr, free := alloc.Calloc(len(buf)+1, 2)
	copy(unsafe.Slice((*uint16)(ptr), len(buf)), buf)
	return (*wchar2_t)(ptr), free
}

func CWLen(s string) int {
	return len(utf16.Encode([]rune(s)))
}

func CWStringCopyTo(dst *wchar2_t, dstSz int, src string) {
	str := unsafe.Slice((*uint16)(unsafe.Pointer(dst)), dstSz)
	if len(src) == 0 {
		str[0] = 0
		return
	}
	n := copy(str[:len(str)-1], utf16.Encode([]rune(src)))
	str[n] = 0
}

func GoWStrSlice(arr **uint16) []string {
	n := alloc.ZeroTermLen(arr)
	return GoWStrSliceN(arr, n)
}

func GoWStrSliceN(arr **uint16, n int) []string {
	if n == 0 {
		return nil
	}
	out := make([]string, 0, n)
	for _, c := range unsafe.Slice(arr, n) {
		out = append(out, GoWString(c))
	}
	return out
}
