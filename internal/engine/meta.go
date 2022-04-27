package engine

import (
	"doraemon/conf"
	"doraemon/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	segMetaFile = "segments.gen" // 存储的元数据文件，包含各种属性信息
)

// Meta 元数据
type Meta struct {
	Version string `json:"version"` // 版本号
	Path    string `json:"path"`    // 存储路径
	// CurrSeg  uint64     `json:"curr_seg"`  // 当前正在使用的seg
	NextSeg  uint64     `json:"next_seg"`  // 下一个segment的命名
	SegCount uint64     `json:"seg_count"` // 当前segment的数量
	SegInfo  []*SegInfo `json:"seg_info"`  // 当前segments的信息

	sync.Mutex
}

// SegInfo 段信息
type SegInfo struct {
	SegID            uint64 `json:"seg_name"`           // 段前缀名
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

// ParseMeta 解析数据
func ParseMeta(c *conf.Config) (*Meta, error) {
	// 文件不存在表示没有相关数据 第一次创建
	segMetaFile = c.Storage.Path + segMetaFile
	if !utils.ExistFile(segMetaFile) {
		log.Debugf("segMetaFile:%s not exist", segMetaFile)
		_, err := os.Create(segMetaFile)
		if err != nil {
			return nil, fmt.Errorf("create segmentsGenFile err: %v", err)
		}
		m := &Meta{
			NextSeg:  0,
			Version:  c.Version,
			Path:     segMetaFile,
			SegCount: 0,
			SegInfo:  nil,
		}
		err = writeSeg(m)
		if err != nil {
			return nil, fmt.Errorf("writeSeg err: %v", err)
		}
		return m, nil
	}

	return readSeg(segMetaFile)
}

// SyncByTicker 定时同步元数据
func (m *Meta) SyncByTicker(ticker *time.Ticker) {
	// 清理计时器
	// defer ticker.Stop()
	for {
		log.Infof("ticker start:%s,seg id :%d", time.Now().Format("2006-01-02 15:04:05"), m.NextSeg)
		err := m.SyncMeta()
		if err != nil {
			log.Errorf("sync meta err:%v", err)
		}
		<-ticker.C
	}
}

// SyncMeta 同步元数据到文件
func (m *Meta) SyncMeta() error {
	err := writeSeg(m)
	if err != nil {
		return fmt.Errorf("writeSeg err: %v", err)
	}
	return nil
}

// UpdateSegMeta 更新段信息
func (m *Meta) UpdateSegMeta(indexCount uint64) error {
	m.Lock()
	defer m.Unlock()

	seg := &SegInfo{
		SegID:   m.NextSeg,
		SegSize: indexCount,
	}
	m.addNewSeg(seg)

	err := m.SyncMeta()
	if err != nil {
		return fmt.Errorf("sync writeSeg err: %v", err)
	}
	return nil
}

// NewSegment 创建新的segment 只创建，更新nextseg，不更新currseg
func (m *Meta) NewSegment() *SegInfo {
	m.Lock()
	defer m.Unlock()
	seg := &SegInfo{
		SegID:   m.NextSeg,
		SegSize: 0,
	}
	m.addNewSeg(seg)
	return seg
}

func (m *Meta) addNewSeg(seg *SegInfo) {
	m.SegInfo = append(m.SegInfo, seg)
	m.SegCount++
	m.NextSeg++

}

func readSeg(segMetaFile string) (*Meta, error) {
	metaByte, err := os.ReadFile(segMetaFile)
	if err != nil {
		return nil, fmt.Errorf("read file err: %v", err)
	}
	h := new(Meta)
	err = json.Unmarshal(metaByte, &h)
	if err != nil {
		return nil, fmt.Errorf("ParseHeader err: %v", err)
	}
	log.Debugf("seg header :%v", h)
	// if h.Path != segMetaFile {
	// 	return nil, fmt.Errorf("segMetaFile:%s path is not equal", segMetaFile)
	// }
	return h, nil
}

func writeSeg(m *Meta) error {
	f, err := os.OpenFile(m.Path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open file err: %v", err)
	}
	defer f.Close()
	b, _ := json.Marshal(m)
	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("write file err: %v", err)
	}
	return nil
}
