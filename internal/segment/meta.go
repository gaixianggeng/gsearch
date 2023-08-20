package segment

import (
	"fmt"
	"sync"
)

// SegMeta 元数据
type SegMeta struct {
	NextSeg  SegID              `json:"next_seg"`  // 下一个segmentid,永远表示下一个新建的segment,seginfos中不存在
	SegCount uint64             `json:"seg_count"` // 当前segment的数量
	SegInfo  map[SegID]*SegInfo `json:"seg_info"`  // 当前segments的信息
	sync.Mutex
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
