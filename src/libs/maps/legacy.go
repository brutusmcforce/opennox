package maps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brutusmcforce/opennox/libs/datapath"
	"github.com/brutusmcforce/opennox/libs/ifs"
)

func NewNativeDownloader(dir string) *NativeDownloader {
	if dir == "" {
		dir = datapath.Maps()
	}
	d := &NativeDownloader{
		dir: dir,
	}
	d.Reset()
	return d
}

type mapDownloadPart struct {
	Ind  uint
	Data []byte
	Next *mapDownloadPart
	Prev *mapDownloadPart
}

type NativeDownloader struct {
	dir        string
	path       string
	downloadOK bool
	file       *atomicFile
	curPart    uint
	sizeCur    uint
	sizeTotal  uint
	sizePend   uint // bytes in out-of-order parts
	first      *mapDownloadPart
	last       *mapDownloadPart
}

func (d *NativeDownloader) WritePart(ind uint, data []byte) {
	if len(data) == 0 || d.file == nil || ind < d.curPart {
		return
	}
	// the parts come from the server: never take more than the size it announced
	if d.sizeCur+d.sizePend+uint(len(data)) > d.sizeTotal {
		return
	}
	if ind == d.curPart {
		d.file.Write(data)
		d.sizeCur += uint(len(data))
		d.curPart++
	} else {
		buf := make([]byte, len(data))
		copy(buf, data)
		v3 := &mapDownloadPart{
			Ind:  ind,
			Data: buf,
			Prev: d.last,
		}
		if d.last != nil {
			d.last.Next = v3
		}
		d.last = v3
		if d.first == nil {
			d.first = v3
		}
		d.sizePend += uint(len(buf))
	}

	// outer loop accounts for packet reordering
	for d.first != nil {
		found := false
		// in-order packets
		for it := d.first; it != nil; it = it.Next {
			if d.curPart != it.Ind {
				continue
			}
			found = true
			d.file.Write(it.Data)
			d.sizeCur += uint(len(it.Data))
			d.sizePend -= uint(len(it.Data))
			d.curPart++
			if v6 := it.Prev; v6 != nil {
				v6.Next = it.Next
			} else {
				d.first = it.Next
			}
			if v7 := it.Next; v7 != nil {
				v7.Prev = it.Prev
			} else {
				d.last = it.Prev
			}
			it.Data = nil
		}
		if !found {
			break
		}
	}
}

// CancelAndCleanup stops the download. The map that was there before stays.
func (d *NativeDownloader) CancelAndCleanup() {
	d.Reset()
}

// CancelAndDelete stops the download. The map that was there before stays.
func (d *NativeDownloader) CancelAndDelete() {
	if d.file != nil {
		d.file.Abort()
	}
	d.file = nil
	d.downloadOK = false
	d.curPart = 1
	d.sizeTotal = 0
	d.sizeCur = 0
	d.sizePend = 0

	var v1 *mapDownloadPart
	for v0 := d.first; v0 != nil; v0 = v1 {
		v1 = v0.Next
		*v0 = mapDownloadPart{}
	}
	d.path = ""
	d.last = nil
	d.first = nil
}

func (d *NativeDownloader) DownloadOK() bool {
	return d.downloadOK
}

func (d *NativeDownloader) Complete() bool {
	return d.downloadOK && d.sizeCur >= d.sizeTotal
}

func (d *NativeDownloader) Progress() float64 {
	return float64(d.sizeCur) / float64(d.sizeTotal)
}

// Finish saves a complete download as the map file.
func (d *NativeDownloader) Finish() error {
	if !d.Complete() || d.file == nil {
		d.Reset()
		return fmt.Errorf("map download incomplete: %d of %d bytes", d.sizeCur, d.sizeTotal)
	}
	f := d.file
	d.file = nil
	d.Reset()
	return f.Commit()
}

// Reset stops the download, if any. The map that was there before stays.
func (d *NativeDownloader) Reset() {
	if d.file != nil {
		d.file.Abort()
	}
	d.file = nil
	d.downloadOK = false
	d.last = nil
	d.first = nil
	d.sizeTotal = 0
	d.sizeCur = 0
	d.sizePend = 0
	d.curPart = 1
	d.path = ""
}

func (d *NativeDownloader) Start(path string, sz uint) error {
	d.Reset()
	if sz > mapFileSizeLimit {
		return fmt.Errorf("map is too large: %d bytes", sz)
	}
	path = ifs.Normalize(path)
	mdir := filepath.Dir(path)
	if _, err := ifs.Stat(mdir); os.IsNotExist(err) {
		ifs.MkdirAll(mdir)
	}
	d.path = path
	f, err := createAtomic(d.path)
	if err != nil {
		return err
	}
	d.file = f
	d.curPart = 1
	d.downloadOK = true
	d.sizeTotal = sz
	return nil
}
