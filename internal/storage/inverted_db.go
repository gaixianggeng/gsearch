package storage

import (
	"bytes"
	"doraemon/pkg/utils"
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

// KvInfo term信息
type KvInfo struct {
	Key   []byte
	Value []byte
}

// TermValues 存储的doc_count、offset、size
type TermValues [3]uint64

// StoragePostings 倒排列表存储到数据库中
func (t *InvertedDB) StoragePostings(token string, values []byte, docCount uint64) error {
	// 写入file，获取写入的size
	size, err := t.storagePostings(values)
	if err != nil {
		return fmt.Errorf("DBUpdatePostings storagePostings err: %v", err)
	}

	// 写入b+tree
	buf := bytes.NewBuffer(nil)
	err = utils.BinaryWrite(buf, docCount)
	if err != nil {
		return fmt.Errorf("BinaryWrite docCount err: %v", err)
	}
	err = utils.BinaryWrite(buf, []uint64{t.offset, size})
	if err != nil {
		return fmt.Errorf("BinaryWrite offset size err: %v", err)
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

// GetTermInfo 获取term关联的倒排地址
func (t *InvertedDB) GetTermInfo(token string) (*TermValues, error) {
	c, err := t.Get([]byte(token))
	if err != nil {
		return nil, fmt.Errorf("GetTermInfo err:%v", err)
	}
	return Bytes2TermVal(c)
}

// GetDocInfo 根据地址获取读取文件
func (t *InvertedDB) GetDocInfo(offset uint64, size uint64) ([]byte, error) {
	page := os.Getpagesize()
	b, err := Mmap(int(t.file.Fd()), int64(offset/uint64(page)), int(offset+size))
	if err != nil {
		return nil, fmt.Errorf("GetDocinfo Mmap err: %v", err)
	}
	return b[offset : offset+size], nil
}

// GetTermCursor 获取遍历游标
func (t *InvertedDB) GetTermCursor(termCh chan KvInfo) {
	t.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(termBucket))
		// b.ForEach(func(k, v []byte) error {
		// 	termCh <- TermInfo{k, v}
		// 	return nil
		// })
		// close(termCh)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			termCh <- KvInfo{k, v}
		}
		close(termCh)
		return nil
	})
}

func (t *InvertedDB) storagePostings(postings []byte) (uint64, error) {
	size, err := t.file.WriteAt(postings, int64(t.offset))
	if err != nil {
		return 0, fmt.Errorf("write storage postings err:%v", err)
	}
	return uint64(size), nil

}

// Close 关闭
func (t *InvertedDB) Close() {
	// stat, _ := t.file.Stat()
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

	log.Debugf("start op bolt:%s", termName)
	db, err := bolt.Open(termName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("file size:%d", stat.Size())
	return &InvertedDB{f, db, uint64(stat.Size())}
}

// Bytes2TermVal 字节转换为TermValues
func Bytes2TermVal(values []byte) (*TermValues, error) {
	var p TermValues
	err := binary.Read(bytes.NewBuffer(values), binary.LittleEndian, &p)
	if err != nil {
		return nil, fmt.Errorf("fetchPostings BinaryRead err: %v", err)
	}

	return &p, nil

}
