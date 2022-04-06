package index

import (
	"brain/internal/storage"
)

// Engine 写入引擎
type Engine struct {
	db *storage.DB

	postingsHashBuf InvertedIndexHash // 倒排索引缓冲区
	bufCount        int64             //倒排索引缓冲区的文档数
	bufSize         int64
	indexCount      int64
	N               int32 // ngram
}

// PostingsList 倒排列表
type PostingsList struct {
	DocID         int64
	positions     []int64
	positionCount int64
	next          *PostingsList
}

//InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	TokenID       int64
	postingList   *PostingsList
	docsCount     int64
	positionCount int64
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[int64]*InvertedIndexValue
