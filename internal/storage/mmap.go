package storage

import "syscall"

// Mmap --
func Mmap(fd int, offset int64, length int) ([]byte, error) {
	return syscall.Mmap(fd, offset, length, syscall.PROT_READ, syscall.MAP_SHARED)
}

// Munmap --
func Munmap(b []byte) error {
	return syscall.Munmap(b)
}
