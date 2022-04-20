package storage

import (
	"brain/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

const termBucket = "term"

// InvertedDB 倒排索引存储库
type InvertedDB struct {
	// tree *bptree.Tree
	file   *os.File
	db     *bolt.DB
	offset uint64
}

// DBUpdatePostings 倒排列表存储到数据库中
func (t *InvertedDB) DBUpdatePostings(token string, values []byte) error {
	// 写入file，获取写入的size
	size, err := t.storagePostings(values)
	if err != nil {
		return fmt.Errorf("DBUpdatePostings storagePostings err: %v", err)
	}

	// 写入b+tree
	buf := bytes.NewBuffer(nil)
	log.Debugf("offset:%d, size:%d", t.offset, size)
	err = utils.BinaryWrite(buf, []uint64{t.offset, size})
	if err != nil {
		return fmt.Errorf("DBUpdatePostings BinaryWrite err: %v", err)
	}

	//update offset
	t.offset += size

	return t.Put([]byte(token), buf.Bytes())
}

// Put 插入term
func (t *InvertedDB) Put(key, value []byte) error {
	return Put(t.db, termBucket, key, value)
}

// Get 通过term获取value
func (t *InvertedDB) Get(key []byte) (value []byte, err error) {
	return Get(t.db, termBucket, key)
}

// GetForwordAddr 获取正排地址
func (t *InvertedDB) GetForwordAddr(token string) (offset uint64, size uint64, err error) {
	c, err := t.Get([]byte(token))
	if err != nil {
		return 0, 0, fmt.Errorf("fetchPostings Get err: %v", err)
	}
	p := make([]uint64, 2)

	err = binary.Read(bytes.NewBuffer(c), binary.LittleEndian, &p)
	if err != nil {
		return 0, 0, fmt.Errorf("fetchPostings BinaryRead err: %v", err)
	}
	return p[0], p[1], nil
}

// GetForwordContent 根据地址获取读取文件
func (t *InvertedDB) GetForwordContent(offset uint64, size uint64) ([]byte, error) {
	page := os.Getpagesize()
	b, err := Mmap(int(t.file.Fd()), int64(offset/uint64(page)), int(offset+size))
	if err != nil {
		return nil, fmt.Errorf("GetForwordContent Mmap err: %v", err)
	}
	return b[offset : offset+size], nil
}

func (t *InvertedDB) storagePostings(postings []byte) (uint64, error) {
	log.Debugf("postings len:%d", len(postings))
	size, err := t.file.WriteAt(postings, int64(t.offset))
	if err != nil {
		return 0, fmt.Errorf("write storage postings err:%v", err)
	}
	return uint64(size), nil

}

// Close 关闭
func (t *InvertedDB) Close() {
	stat, _ := t.file.Stat()
	log.Debugf("close,offset:%d,file size:%d", t.offset, stat.Size())
	t.file.Close()
	t.db.Close()
}

// NewInvertedDB 初始化
func NewInvertedDB(termName, postingsName string) *InvertedDB {
	// tree, err := bptree.NewTree(bTreeName)
	// if err != nil {
	// 	panic(err)
	// }
	f, err := os.OpenFile(postingsName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(termName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("file size:%d", stat.Size())
	return &InvertedDB{f, db, uint64(stat.Size())}
}
