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
}

// InvertedItem 写入库文件结构
type InvertedItem struct {
	PostingsList []byte
	PostingsLen  uint64
	DocCount     uint64
}

// DBUpdatePostings 倒排列表存储到数据库中
func (t *InvertedDB) DBUpdatePostings(tokenID uint64, values *InvertedItem) error {
	// 写入file
	// t.file.Write( )

	// 获取file的offset

	// 写入b+tree

	cByte, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("DBUpdatePostings json.Marshal err: %v", err)
	}
	return t.tree.Insert(tokenID, cByte)
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
	return &InvertedDB{tree, f}
}
