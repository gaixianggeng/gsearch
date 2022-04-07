package storage

import (
	"brain/internal/storage/bptree"
	"bytes"
)

// InvertedDB 倒排索引存储库
type InvertedDB struct {
	client *bptree.Tree
}

// DBUpdatePostings 倒排列表存储到数据库中
func DBUpdatePostings(db *InvertedDB, tokenID uint64, docsCount uint64, buf *bytes.Buffer, bufSize uint64) error {
	db.client.Insert(tokenID, "")
	return nil
}

// NewInvertedDB 初始化
func NewInvertedDB() *InvertedDB {
	tree, err := bptree.NewTree("./inverted.db")
	if err != nil {
		panic(err)
	}
	return &InvertedDB{client: tree}
}
