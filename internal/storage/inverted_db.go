package storage

import (
	"brain/internal/storage/bptree"
	"encoding/json"
	"fmt"
	"os"
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
	// t.file.Write(values)

	// 获取file的offset

	// 写入b+tree

	cByte, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("DBUpdatePostings json.Marshal err: %v", err)
	}
	return t.tree.Insert(tokenID, cByte)
}

func (t *InvertedDB) storagePostings(postings []byte) {

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
