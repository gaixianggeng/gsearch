package index

import (
	"brain/internal/storage"
)

// Engine 写入引擎
type Engine struct {
	db *storage.DB

	postingsHashBuf InvertedIndexHash // 倒排索引缓冲区
	bufCount        uint64            //倒排索引缓冲区的文档数
	bufSize         uint64
	indexCount      uint64
	N               int32 // ngram
}

// PostingsList 倒排列表
type PostingsList struct {
	DocID         uint64
	positions     []uint64
	positionCount uint64
	next          *PostingsList
}

//InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	TokenID       uint64
	postingList   *PostingsList
	docsCount     uint64
	positionCount uint64
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[uint64]*InvertedIndexValue
