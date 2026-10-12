package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/cnxz"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"os"
)

var worldMapSaveService = worldMapSave

func worldSaveDrawables() int {
	cli := GetClient().Cli()
	srv := GetServer().S()
	for d := cli.Objs.FirstList1(); d != nil; d = d.Next() {
		typ := int(d.TypeIDVal)
		if Sub_4E3AD0(typ) != 0 || !Sub_4E3B80(typ) {
			continue
		}
		u := srv.NewObjectByTypeID(cli.Things.TypeByInd(typ).ID())
		originalCode := u.NetCode
		u.PosVec.X = float32(float64(d.PosVec.X) + 0.5)
		u.PosVec.Y = float32(float64(d.PosVec.Y) + 0.5)
		u.NetCode = d.NetCode32
		u.Extent = d.NetCode32
		u.ScriptIDVal = int(d.NetCode32)
		u.ObjFlags = d.ObjFlags
		u.Field5 = d.Flags70Val
		Nox_xxx_xfer_saveObj51DF90(cryptfile.Global(), u)
		u.NetCode = originalCode
		srv.Objs.FreeObject(u)
	}
	return 1
}
func worldMapSave(path string, compression int) bool {
	// The former fixed-size C buffers required a four-byte extension and <1024
	// bytes. Reject inputs outside that domain before deriving the output path.
	if len(path) < 4 || len(path) >= 1024 {
		return false
	}
	compressed := path[:len(path)-4] + alloc.GoString(memmap.PtrUint8(0x587000, 253112))
	if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, 19); err != nil {
		if !os.IsNotExist(err) {
			binfile.Log.Println(err)
		}
		return false
	}
	cf := cryptfile.Global()
	cf.WriteU32(0xfadeface)
	position := cf.Flush()
	if !worldMapWrite() {
		cryptfile.Close()
		return false
	}
	cryptfile.Global().WriteChecksumAt(int64(position))
	cryptfile.Close()
	if compression != 0 {
		if err := cnxz.CompressFile(path, compressed); err != nil {
			mapLog.Println(err)
			return false
		}
	}
	return true
}
func worldMapWrite() bool {
	x, y := memmap.PtrInt32(0x5D4594, 2487252), memmap.PtrInt32(0x5D4594, 2487256)
	*x, *y = 256, 256
	GetServer().S().Walls.EachWallXxx(func(w *server.Wall) bool {
		if int32(w.X5) < *x {
			*x = int32(w.X5)
		}
		if int32(w.Y6) < *y {
			*y = int32(w.Y6)
		}
		return true
	})
	cf := cryptfile.Global()
	cf.WriteU32(uint32(*x))
	cf.WriteU32(uint32(*y))
	*memmap.PtrInt32(0x5D4594, 739980) = *x
	*memmap.PtrInt32(0x5D4594, 739984) = *y
	Nox_xxx_mapSetWallInGlobalDir0pr1_5004D0()
	if Nox_xxx_mapWriteSectionsMB_426E20(nil) == 0 {
		cryptfile.Close()
		return false
	}
	Nox_xxx_map_5004F0()
	cryptfile.Global().WriteU8(0)
	return true
}
