package segment

import (
	"gsearch/internal/storage"
)

var (
	// TermDBSuffix term db suffix
	TermDBSuffix = ".term"
	// InvertedDBSuffix inverted db suffix
	InvertedDBSuffix = ".inverted"
	// ForwardDBSuffix forward db suffix
	ForwardDBSuffix = ".forward"
)
var (
	termName     = ""
	invertedName = ""
	forwardName  = ""
)

// InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token            string             // 词元
	PostingsList     *PostingsList      // 文档编号的序列
	DocCount         uint64             // 词元关联的文档数量
	DocPositionCount uint64             // 词元在所有文档中出现的次数 查询使用,用于计算相关性，写入的时候暂时用不到
	TermValues       *storage.TermValue // 存储的doc_count、offset、size
}

// LoserTree 败者树 用于多路归并
type LoserTree struct {
	tree     []int // 索引表示顺序，0表示最小值，value表示对应的leaves的index
	leaves   []*TermNode
	leavesCh []chan storage.KvInfo
}

// TermNode 词元节点
type TermNode struct {
	*storage.KvInfo
	Seg *Segment // 主要用来调用Inverted相关方法
}

// Mode 查询or索引模式
type Mode int32

const (
	// SearchMode 查询模式
	SearchMode Mode = 1
	// IndexMode 索引模式
	IndexMode Mode = 2
	// MergeMode seg merge模式
	MergeMode Mode = 3
)

// SegID --
type SegID int64

// SegInfo 段信息
type SegInfo struct {
	SegID            SegID  `json:"seg_name"`           // 段前缀名
	SegSize          uint64 `json:"seg_size"`           // 写入doc数量
	InvertedFileSize uint64 `json:"inverted_file_size"` // 写入inverted文件大小
	ForwardFileSize  uint64 `json:"forward_file_size"`  // 写入forward文件大小
	DelSize          uint64 `json:"del_size"`           // 删除文档数量
	DelFileSize      uint64 `json:"del_file_size"`      // 删除文档文件大小
	TermSize         uint64 `json:"term_size"`          // term文档文件大小
	TermFileSize     uint64 `json:"term_file_size"`     // term文件大小
	ReferenceCount   uint64 `json:"reference_count"`    // 引用计数
	IsReading        bool   `json:"is_reading"`         // 是否正在被读取
	IsMerging        bool   `json:"is_merging"`         // 是否正在参与合并
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*InvertedIndexValue

// PostingsList 倒排列表
type PostingsList struct {
	DocID         uint64
	Positions     []uint64
	PositionCount uint64
	Next          *PostingsList
}
