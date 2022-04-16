package storage

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

// InvertedDB 倒排索引存储库
type InvertedDB struct {
	// tree *bptree.Tree
	file   *os.File
	db     *bolt.DB
	offset uint64
}

// DBUpdatePostings 倒排列表存储到数据库中
func (t *InvertedDB) DBUpdatePostings(token string, values []byte) error {
	// 写入file
	size, err := t.storagePostings(values)
	if err != nil {
		return fmt.Errorf("DBUpdatePostings storagePostings err: %v", err)
	}

	return nil

	// 获取file的offset

	// 写入b+tree
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, t.offset)
	log.Debug(string(value))

	//update offset
	t.offset += size

	return t.Put([]byte(token), value)
}

// Put 插入term
func (t *InvertedDB) Put(key, value []byte) error {

	return nil
}

// Get 通过term获取value
func (t *InvertedDB) Get(key []byte) (value []byte, err error) {

	return nil, nil
}

func (t *InvertedDB) storagePostings(postings []byte) (uint64, error) {
	log.Debugf("postings len:%d", len(postings))
	size, err := t.file.Write(postings)
	if err != nil {
		return 0, fmt.Errorf("write storage postings err:%v", err)
	}
	return uint64(size), nil

}

// Close 关闭
func (t *InvertedDB) Close() {
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

	return &InvertedDB{f, db, uint64(stat.Size())}
}
