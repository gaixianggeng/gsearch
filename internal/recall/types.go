package recall

import "gsearch/internal/segment"

// 用于实现排序map
type queryTokenHash struct {
	token         string
	invertedIndex *segment.InvertedIndexValue
	fetchPostings *segment.PostingsList
}

// SearchItem 查询结果
type SearchItem struct {
	DocID uint64
	Score float64
}

// Recalls 召回结果
type Recalls []*SearchItem

// token游标 标识当前位置
type searchCursor struct {
	doc     *segment.PostingsList // 文档编号的序列
	current *segment.PostingsList // 当前的文档编号
}

// 短语游标
type phraseCursor struct {
	positions []uint64 // 位置信息
	base      uint64   // 词元在查询中的位置
	current   *uint64  // 当前的位置信息
	index     uint     // 当前的位置index
}
