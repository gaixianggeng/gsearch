package storage

import (
	"fmt"
	"gsearch/pkg/utils/log"
	"os"
	"syscall"
)

const maxMapSize = 0x1000000000 // 64GB
const maxMmapStep = 1 << 30     // 1GB

// OpenFile --
func OpenFile() {

}

func writeMMap() {
	file, err := os.OpenFile("my.db", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := os.Stat("my.db")
	if err != nil {
		panic(err)
	}

	size, err := mmapSize(int(stat.Size()))
	if err != nil {
		panic(err)
	}
	syscall.Ftruncate(int(file.Fd()), int64(size))

	b, err := syscall.Mmap(
		int(file.Fd()),
		0,
		size,
		syscall.PROT_WRITE|syscall.PROT_READ,
		syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	for index, bb := range []byte("Hello world") {
		b[index] = bb
	}

	err = syscall.Munmap(b)
	if err != nil {
		panic(err)
	}
}

func mmapSize(size int) (int, error) {
	// Double the size from 32KB until 1GB.
	// 最小每次增加32kb
	// 如果size小于1GB，每次double递增
	for i := uint(15); i <= 30; i++ {
		if size <= 1<<i {
			return 1 << i, nil
		}
	}

	// Verify the requested size is not above the maximum allowed.
	if size > maxMapSize {
		return 0, fmt.Errorf("mmap too large")
	}

	// If larger than 1GB then grow by 1GB at a time.
	// 大于1GB，每次增加1GB,第一次补全1GB
	sz := int64(size)
	if remainder := sz % int64(maxMmapStep); remainder > 0 {
		sz += int64(maxMmapStep) - remainder
	}

	// Ensure that the mmap size is a multiple of the page size.
	// This should always be true since we're incrementing in MBs.
	// 获取pagesize
	pageSize := int64(os.Getpagesize())
	if (sz % pageSize) != 0 {
		sz = ((sz / pageSize) + 1) * pageSize
	}

	// If we've exceeded the max size then only grow up to the max size.
	if sz > maxMapSize {
		sz = maxMapSize
	}

	return int(sz), nil
}

func read() {
	file, err := os.OpenFile("my.db", os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := os.Stat("my.db")
	if err != nil {
		panic(err)
	}
	log.Debug(stat.Size())

	// b, err := syscall.Mmap(int(file.Fd()), 0, int(stat.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	log.Debug(int(file.Fd()))
	log.Debug(file.Fd())
	b, err := syscall.Mmap(int(file.Fd()), 0, 4, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(b)

	log.Debug(string(b))
	log.Debug(len(b))
}
