package alloc

import (
	"bytes"
	"unsafe"
)

func Memset(ptr unsafe.Pointer, v byte, size uintptr) unsafe.Pointer {
	logMemWrite(ptr, size)
	buf := unsafe.Slice((*byte)(ptr), int(size))
	if v == 0 {
		clear(buf)
	} else if len(buf) != 0 {
		buf[0] = v
		for n := 1; n < len(buf); {
			n += copy(buf[n:], buf[:n])
		}
	}
	return ptr
}

func Memcpy(dst, src unsafe.Pointer, size uintptr) unsafe.Pointer {
	logMemRead(src, size)
	logMemWrite(dst, size)
	copy(unsafe.Slice((*byte)(dst), int(size)), unsafe.Slice((*byte)(src), int(size)))
	return dst
}

func Memcmp(ptr1, ptr2 unsafe.Pointer, size uintptr) int {
	logMemRead(ptr1, size)
	logMemRead(ptr2, size)
	return bytes.Compare(unsafe.Slice((*byte)(ptr1), int(size)), unsafe.Slice((*byte)(ptr2), int(size)))
}

func Strcpy(dst, src unsafe.Pointer) unsafe.Pointer {
	n := uintptr(strlen((*byte)(src)))
	logMemReadString(src, n+1)
	logMemWriteString(dst, n+1)
	copy(unsafe.Slice((*byte)(dst), int(n+1)), unsafe.Slice((*byte)(src), int(n+1)))
	return dst
}

func Strcat(dst, src unsafe.Pointer) unsafe.Pointer {
	ns := uintptr(strlen((*byte)(src)))
	nd := uintptr(strlen((*byte)(dst)))
	logMemReadString(src, ns+1)
	logMemWriteString(dst, nd+ns+1)
	copy(unsafe.Slice((*byte)(unsafe.Add(dst, nd)), int(ns+1)), unsafe.Slice((*byte)(src), int(ns+1)))
	return dst
}

func Strcmp(str1, str2 unsafe.Pointer) int {
	n1 := uintptr(strlen((*byte)(str1)))
	n2 := uintptr(strlen((*byte)(str2)))
	logMemReadString(str1, n1+1)
	logMemReadString(str2, n2+1)
	return bytes.Compare(unsafe.Slice((*byte)(str1), int(n1)), unsafe.Slice((*byte)(str2), int(n2)))
}
