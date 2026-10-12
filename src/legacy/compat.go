package legacy

import (
	"time"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/env"
	"github.com/brutusmcforce/opennox/libs/platform"
)

type noxSystemTime struct {
	wYear, wMonth, wDayOfWeek, wDay        uint16
	wHour, wMinute, wSecond, wMilliseconds uint16
}

var _ = [1]struct{}{}[16-unsafe.Sizeof(noxSystemTime{})]

func noxGetLocalTime(p *noxSystemTime) {
	tm := time.Now()
	if env.IsE2E() {
		tm = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC).Add(platform.Ticks())
	}
	p.wYear = uint16(tm.Year())
	p.wMonth = uint16(tm.Month())
	p.wDayOfWeek = uint16(tm.Weekday())
	p.wDay = uint16(tm.Day())
	p.wHour = uint16(tm.Hour())
	p.wMinute = uint16(tm.Minute())
	p.wSecond = uint16(tm.Second())
	p.wMilliseconds = uint16(tm.Nanosecond() / 1e6)
}
