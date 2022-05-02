package segment

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

// NewSegment 创建新的segment 只创建，更新nextseg，不更新currseg
func NewSegment(segID SegID) *SegInfo {
	seg := &SegInfo{
		SegID:   segID,
		SegSize: 0,
	}
	return seg
}
