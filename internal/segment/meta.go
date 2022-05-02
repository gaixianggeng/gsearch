package segment

import (
	"fmt"
	"sync"
)

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

// SegMeta 元数据
type SegMeta struct {
	NextSeg  SegID              `json:"next_seg"`  // 下一个segmentid,永远表示下一个新建的segment,seginfos中不存在
	SegCount uint64             `json:"seg_count"` // 当前segment的数量
	SegInfo  map[SegID]*SegInfo `json:"seg_info"`  // 当前segments的信息

	sync.Mutex
}

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

// newSegment 创建新的segment 只创建，更新nextseg，不更新currseg
func newSegmentInfo(segID SegID) *SegInfo {
	seg := &SegInfo{
		SegID:   segID,
		SegSize: 0,
	}
	return seg
}

// UpdateSegMeta 更新段信息
func (m *SegMeta) UpdateSegMeta(segID SegID, indexCount uint64) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.SegInfo[segID]; !ok {
		return fmt.Errorf("seg:%d is not exist", segID)
	}
	m.SegInfo[segID].SegSize = indexCount
	return nil
}

// NewSegmentItem 创建新的segment 只创建，更新nextseg，不更新currseg
func (m *SegMeta) NewSegmentItem() error {
	m.Lock()
	defer m.Unlock()
	seg := newSegmentInfo(m.NextSeg)
	if _, ok := m.SegInfo[SegID(seg.SegID)]; ok {
		return fmt.Errorf("seg:%d is exist", seg.SegID)
	}
	m.SegInfo[SegID(seg.SegID)] = seg
	m.SegCount++
	m.NextSeg++
	return nil
}
