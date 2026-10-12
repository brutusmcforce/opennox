package opennox

import (
	"github.com/brutusmcforce/opennox/libs/strman"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

var (
	strMan     = strman.New()
	strManDone = false
)

func StrmanReadFile(path string) error {
	if strManDone {
		return nil
	}
	err := strMan.ReadFile(path)
	if err != nil {
		return err
	}
	strManDone = true
	return nil
}

func Nox_strman_free_410020() {
	strMan = strman.New()
	strManDone = false
	alloc.FreeStrings()
}
