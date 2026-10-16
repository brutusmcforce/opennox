package sdl

import (
	"sync"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"

	"github.com/brutusmcforce/opennox/libs/dynlib"
)

// The system file dialogs: xdg-desktop-portal on Linux (falling back to zenity, or kdialog), NSOpenPanel on macOS,
// IFileDialog on Windows.

const fileDialogOpenFolder = 2 // SDL_FILEDIALOG_OPENFOLDER

var dialogFn struct {
	once     sync.Once
	err      error
	show     uintptr
	newProps uintptr
	setStr   uintptr
	delProps uintptr
	pump     uintptr
	quitSub  uintptr
	callback uintptr // the SDL_DialogFileCallback, created once: purego can't free callbacks
}

type dialogResult struct {
	path string
	err  error
}

// dialogDone receives the result from the callback, which SDL may call on another thread.
var dialogDone = make(chan dialogResult, 1)

func loadDialog() error {
	dialogFn.once.Do(func() {
		if dialogFn.err = Load(); dialogFn.err != nil {
			return
		}
		dialogFn.err = loadedLib.Syms(map[string]*uintptr{
			"SDL_ShowFileDialogWithProperties": &dialogFn.show,
			"SDL_CreateProperties":             &dialogFn.newProps,
			"SDL_SetStringProperty":            &dialogFn.setStr,
			"SDL_DestroyProperties":            &dialogFn.delProps,
			"SDL_PumpEvents":                   &dialogFn.pump,
			"SDL_QuitSubSystem":                &dialogFn.quitSub,
		})
		if dialogFn.err != nil {
			return
		}
		// void (*)(void *userdata, const char * const *filelist, int filter)
		dialogFn.callback = purego.NewCallback(func(userdata, filelist, filter uintptr) uintptr {
			var res dialogResult
			switch {
			case filelist == 0:
				res.err = GetError()
			case *(*uintptr)(unsafe.Pointer(filelist)) != 0: //nolint:govet // native memory
				// the strings are only valid during the callback
				res.path = dynlib.GoString(*(*uintptr)(unsafe.Pointer(filelist))) //nolint:govet // native memory
			}
			// filelist[0] == NULL: cancelled, an empty path
			dialogDone <- res
			return 0
		})
	})
	return dialogFn.err
}

// ShowOpenFolderDialog asks the user for a folder with the system's dialog and waits for the answer. It returns ""
// if the user cancelled. It works without a window (and before one is created), but must be called from the main thread.
func ShowOpenFolderDialog(title string) (string, error) {
	if err := loadDialog(); err != nil {
		return "", err
	}
	// the dialogs need the video subsystem (and its event loop) on some platforms
	if err := InitSubSystem(INIT_VIDEO); err != nil {
		return "", err
	}
	defer call(dialogFn.quitSub, INIT_VIDEO)
	props := uint32(call(dialogFn.newProps))
	if props == 0 {
		return "", GetError()
	}
	defer call(dialogFn.delProps, uintptr(props))
	key, val := dynlib.CString("SDL.filedialog.title"), dynlib.CString(title)
	call(dialogFn.setStr, uintptr(props), uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(val)))
	call(dialogFn.show, fileDialogOpenFolder, dialogFn.callback, 0, uintptr(props))
	for {
		select {
		case res := <-dialogDone:
			return res.path, res.err
		default:
		}
		call(dialogFn.pump)
		time.Sleep(10 * time.Millisecond)
	}
}
