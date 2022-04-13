package index

import (
	"brain/internal/storage"
)

// Engine 写入引擎
type Engine struct {
	forwardDB  *storage.ForwardDB
	invertedDB *storage.InvertedDB
	tokenDB    *storage.TokenDB

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
	Token         string
	postingsList  *PostingsList
	docsCount     uint64
	positionCount uint64 // 查询使用，写入的时候暂时用不到
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*InvertedIndexValue
