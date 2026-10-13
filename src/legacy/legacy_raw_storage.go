package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

// legacyGlobalStorage keeps the original field widths and dimensions in unmanaged C-heap memory.
// alloc.New uses calloc; this package-lifetime allocation is intentionally never freed.
type legacyGlobalStorage struct {
	nox_xxx_aClosewoodengat_587000_133480 unsafe.Pointer
	dword_5d4594_831236                   uint32
	dword_5d4594_1309720                  uint32
	dword_587000_127004                   unsafe.Pointer
	nox_alloc_tradeItems_2386496          unsafe.Pointer
	dword_587000_93164                    unsafe.Pointer
	nox_alloc_screenParticles_806044      unsafe.Pointer
	dword_587000_122852                   unsafe.Pointer
	dword_587000_81128                    unsafe.Pointer
	dword_5d4594_1090048                  *gui.Window
	dword_5d4594_1090100                  *gui.Window
	nox_alloc_monsterList_2386220         unsafe.Pointer
	nox_alloc_tradeSession_2386492        unsafe.Pointer
	nox_alloc_spawn_2386216               unsafe.Pointer
	nox_alloc_magicEnt_1569668            unsafe.Pointer
	dword_5d4594_1548532                  unsafe.Pointer
	array_5D4594_1049872                  [9]uint32
	nox_screenParticles_head              *Nox_screenParticle
	dword_5d4594_806052                   *Nox_screenParticle
	dword_5d4594_805984                   unsafe.Pointer
	nox_win_unk1                          *gui.Window
	nox_draw_curDrawData_3799572          *noxrender.RenderData
	dword_5d4594_830236                   unsafe.Pointer
	dword_5d4594_830232                   unsafe.Pointer
	nox_win_unk5                          *gui.Window
	dword_5d4594_1062452                  *gui.Window
	nox_client_inventory_grid_1050020     [84]uiInventoryCell
	nox_pixbuffer_rows_3798784            **byte
	nox_video_tileBuf_ptr_3798796         unsafe.Pointer
	nox_video_tileBuf_end_3798844         unsafe.Pointer
	nox_server_gameSettingsUpdated        int32
	nox_wnd_xxx_1309740                   *gui.Anim
	dword_5d4594_1522616                  *gui.Window
	dword_5d4594_1522620                  *gui.Window
	dword_5d4594_1522624                  *gui.Window
	dword_5d4594_1522628                  *gui.Window
	nox_wnd_xxx_1522608                   *gui.Anim
	nox_gui_itemAmount_item_1319256       unsafe.Pointer
	nox_gui_itemAmount_dialog_1319228     unsafe.Pointer
	dword_5d4594_1321236                  *gui.Window
	dword_5d4594_1321240                  *gui.Window
	dword_5d4594_1321244                  *gui.Window
	dword_5d4594_1321248                  *gui.Window
	nox_common_maplist                    legacyListNode
}

var legacyGlobals, _ = alloc.New(legacyGlobalStorage{})

var byte_581450, _ = alloc.Make([]byte(nil), 23472)
var byte_587000, _ = alloc.Make([]byte(nil), 316820)
var byte_5D4594, _ = alloc.Make([]byte(nil), 2598284)
var byte_973CE0, _ = alloc.Make([]byte(nil), 568)
var byte_973F18, _ = alloc.Make([]byte(nil), 44881)
var byte_85B3FC, _ = alloc.Make([]byte(nil), 1029636)
var byte_852978, _ = alloc.Make([]byte(nil), 40)
var byte_973A20, _ = alloc.Make([]byte(nil), 704)
