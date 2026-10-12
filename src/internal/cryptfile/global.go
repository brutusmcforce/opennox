package cryptfile

var global *CryptFile

// Hooks for a write buffer outside of this package (the legacy C code buffers its small writes, see legacy/crypt_wbuf.c).
var (
	// Sync is called before every operation on the global file, to write the buffered data first.
	Sync func()
	// OnGlobal is called when the global file changes (nil when it's closed).
	OnGlobal func(cf *CryptFile)
)

// sync writes the data buffered outside of this package, if the file is the global one.
func (f *CryptFile) sync() {
	if f != nil && f == global && Sync != nil {
		Sync()
	}
}

func Global() *CryptFile {
	return global
}

func SetGlobal(cf *CryptFile) {
	if global != nil && global != cf {
		global.sync()
	}
	global = cf
	if OnGlobal != nil {
		OnGlobal(cf)
	}
}

func OpenGlobal(path string, cmode Mode, key int) error {
	Close()
	cf, err := OpenFile(path, cmode, key)
	if err != nil {
		return err
	}
	SetGlobal(cf)
	return nil
}

func Close() error {
	if global == nil {
		return nil
	}
	err := global.Close() // writes the buffered data
	SetGlobal(nil)
	return err
}
