package storage

import (
	"brain/internal/storage/bptree"
	"encoding/binary"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// InvertedDB 倒排索引存储库
type InvertedDB struct {
	tree *bptree.Tree
	file *os.File

	offset uint64
}

// DBUpdatePostings 倒排列表存储到数据库中
func (t *InvertedDB) DBUpdatePostings(tokenID uint64, values []byte) error {
	// 写入file
	size, err := t.storagePostings(values)
	if err != nil {
		return fmt.Errorf("DBUpdatePostings storagePostings err: %v", err)
	}

	// 获取file的offset

	// 写入b+tree
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, t.offset)
	log.Debug(string(value))

	//update offset
	t.offset += size

	return t.tree.Insert(tokenID, value)
}

func (t *InvertedDB) storagePostings(postings []byte) (uint64, error) {
	size, err := t.file.Write(postings)
	if err != nil {
		return 0, fmt.Errorf("write storage postings err:%v", err)
	}
	return uint64(size), nil

}

// NewInvertedDB 初始化
func NewInvertedDB(bTreeName, postingsName string) *InvertedDB {
	tree, err := bptree.NewTree(bTreeName)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(postingsName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return &InvertedDB{tree, f, uint64(stat.Size())}
}
