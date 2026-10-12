package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func internCStr(s string) *int8 {
	p := alloc.InternCString(s)
	return (*int8)(unsafe.Pointer(p))
}

func internWStr(s string) *wchar2_t {
	p := alloc.InternCString16(s)
	return (*wchar2_t)(unsafe.Pointer(p))
}

func nox_strman_loadString_40F1D0(name *int8, strOut **int8, srcFile *int8, srcLine int) *wchar2_t {
	if strOut != nil {
		*strOut = nil
		v, _ := GetServer().S().Strings().GetVariantInFile(strman.ID(GoStringP(unsafe.Pointer(name))), GoStringP(unsafe.Pointer(srcFile)))
		*strOut = internCStr(v.Str2)
		return internWStr(v.Str)
	}
	s := GetServer().S().Strings().GetStringInFile(strman.ID(GoStringP(unsafe.Pointer(name))), GoStringP(unsafe.Pointer(srcFile)))
	return internWStr(s)
}

func nox_strman_get_lang_code() int {
	return GetServer().S().Strings().Lang()
}
