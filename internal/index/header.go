package index

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	segmentsGenFile = "segments.gen" // 存储的元数据文件，包含各种属性信息
)

// Header 头信息
type Header struct {
	Version  float64    `json:"version"`   // 版本号
	NextSeg  int64      `json:"next_seg"`  // 下一个segment的命名
	SegCount int64      `json:"seg_count"` // 当前segment的数量
	SegInfo  []*segInfo `json:"seg_info"`  // 当前segments的信息
}

// segInfo 段信息
type segInfo struct {
	SegName string `json:"seg_name"` // 段前缀名

	SegSize          int64 `json:"seg_size"`           // 写入doc数量
	InvertedFileSize int64 `json:"inverted_file_size"` // 写入inverted文件大小
	ForwardFileSize  int64 `json:"forward_file_size"`  // 写入forward文件大小
	DelSize          int64 `json:"del_size"`           // 删除文档数量
	DelFileSize      int64 `json:"del_file_size"`      // 删除文档文件大小
	TermSize         int64 `json:"term_size"`          // term文档文件大小
	TermFileSize     int64 `json:"term_file_size"`     // term文件大小
}

// ParseHeader 解析数据
func ParseHeader() *Header {
	f, err := os.ReadFile(segmentsGenFile)
	if err != nil {
		panic(err)
	}
	log.Debugf("seg:%s", f)
	return nil
}
