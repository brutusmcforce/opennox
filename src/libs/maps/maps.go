package maps

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/brutusmcforce/opennox/libs/ifs"
)

const (
	Dir           = "maps"
	Ext           = ".map"
	SoloPrefixWar = "war"
	SoloPrefixWiz = "wiz"
	SoloPrefixCon = "con"
)

type MapItem struct {
	Field6 int
	Info
}

type MapList []Info

func (m *MapList) Sort() {
	list := *m
	sort.Slice(list, func(i, j int) bool {
		return list[i].Filename < list[j].Filename
	})
	*m = list
}

type MapItemList []MapItem

func (m *MapItemList) Sort() {
	list := *m
	sort.Slice(list, func(i, j int) bool {
		return list[i].Filename < list[j].Filename
	})
	*m = list
}

type ScanOptions struct {
	Solo bool // don't skip solo maps
}

func Scan(log *slog.Logger, path string, opts *ScanOptions) (MapList, error) {
	if opts == nil {
		opts = &ScanOptions{}
	}
	dir, err := ifs.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	var (
		list MapList
		last error
	)
	for {
		page, err := dir.ReadDir(100)
		if err == io.EOF {
			break
		} else if err != nil {
			return list, err
		}
		for _, fi := range page {
			if !fi.IsDir() {
				continue
			}
			name := fi.Name()
			lname := strings.ToLower(name)
			if !opts.Solo && IsSoloName(lname) {
				continue
			}
			info, err := ReadMapInfo(filepath.Join(path, name))
			if os.IsNotExist(err) {
				continue
			} else if err != nil {
				log.Warn("invalid map file", "name", name, "err", err)
				last = err
				continue
			}
			list = append(list, *info)
		}
	}
	list.Sort()
	return list, last
}

// IsSoloName reports whether a map name is one of the solo campaign maps: con, war or wiz,
// a two-digit chapter and a letter (e.g. con01a). Other maps may start with these letters too.
func IsSoloName(name string) bool {
	name = strings.ToLower(name)
	if len(name) != 6 {
		return false
	}
	switch name[:3] {
	case SoloPrefixCon, SoloPrefixWar, SoloPrefixWiz:
	default:
		return false
	}
	isDigit := func(c byte) bool { return c >= '0' && c <= '9' }
	return isDigit(name[3]) && isDigit(name[4]) && name[5] >= 'a' && name[5] <= 'z'
}

func ReadMapInfo(dir string) (*Info, error) {
	name := filepath.Base(dir)
	path := filepath.Join(dir, name+Ext)
	rc, err := ifs.Open(path)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	fi, err := rc.Stat()
	if err != nil {
		return nil, err
	}

	r, err := NewReader(rc)
	if err != nil {
		return nil, err
	}

	info, err := r.ReadInfo()
	if info != nil {
		info.Filename = name
		info.Size = int(fi.Size())
	}
	return info, err
}

func ReadMap(dir string) (*Map, error) {
	name := filepath.Base(dir)
	path := filepath.Join(dir, name+Ext)
	rc, err := ifs.Open(path)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	fi, err := rc.Stat()
	if err != nil {
		return nil, err
	}

	r, err := NewReader(rc)
	if err != nil {
		return nil, err
	}
	err = r.ReadSections()
	m := r.Map()
	if m != nil {
		m.Info.Filename = name
		m.Info.Size = int(fi.Size())
	}
	return m, err
}
