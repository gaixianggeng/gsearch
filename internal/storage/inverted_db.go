package storage

import (
	"brain/internal/storage/bptree"
	"encoding/json"
	"fmt"
)

// InvertedDB 倒排索引存储库
type InvertedDB struct {
	client *bptree.Tree
}

// InvertedItem 写入库文件结构
type InvertedItem struct {
	PostingsList []byte
	PostingsSize uint64
	DocCount     uint64
}

// DBUpdatePostings 倒排列表存储到数据库中
func DBUpdatePostings(
	db *InvertedDB, tokenID uint64, content *InvertedItem) error {
	cByte, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("DBUpdatePostings json.Marshal err: %v", err)
	}
	return db.client.Insert(tokenID, cByte)
}

// NewInvertedDB 初始化
func NewInvertedDB() *InvertedDB {
	tree, err := bptree.NewTree("./inverted.db")
	if err != nil {
		panic(err)
	}
	return &InvertedDB{client: tree}
}
