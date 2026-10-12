package alloc

import "syscall"

func rawMap(size int) []byte {
	data, err := syscall.Mmap(-1, 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		return nil
	}
	return data
}
func rawUnmap(data []byte) error { return syscall.Munmap(data) }
