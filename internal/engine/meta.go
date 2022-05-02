package engine

import (
	"doraemon/conf"
	"doraemon/internal/segment"
	"doraemon/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	metaFile = "segments.json" // 存储的元数据文件，包含各种属性信息
)

// Meta 元数据
type Meta struct {
	sync.RWMutex
	Version    string           `json:"version"` // 版本号
	IndexCount uint64           `json:"index"`
	SegMeta    *segment.SegMeta `json:"seg_meta"`
	path       string           `json:"-"` // 元数据文件路径
}

// ParseMeta 解析数据
func ParseMeta(c *conf.Config) (*Meta, error) {
	// 文件不存在表示没有相关数据 第一次创建
	metaFile = c.Storage.Path + metaFile
	if !utils.ExistFile(metaFile) {
		log.Debugf("segMetaFile:%s not exist", metaFile)
		_, err := os.Create(metaFile)
		if err != nil {
			return nil, fmt.Errorf("create segmentsGenFile err: %v", err)
		}
		m := &Meta{
			Version: c.Version,
			path:    metaFile,
			SegMeta: &segment.SegMeta{
				NextSeg:  0,
				SegCount: 0,
				SegInfo:  make(map[segment.SegID]*segment.SegInfo, 0),
			},
			//TODO: 初始化读取正排数据
			IndexCount: 0,
		}
		err = writeMeta(m)
		if err != nil {
			return nil, fmt.Errorf("writeSeg err: %v", err)
		}
		return m, nil
	}

	return readMeta(metaFile)
}

// SyncByTicker 定时同步元数据
func (m *Meta) SyncByTicker(ticker *time.Ticker) {
	// 清理计时器
	// defer ticker.Stop()
	for {
		log.Infof("ticker start:%s,next seg id :%d",
			time.Now().Format("2006-01-02 15:04:05"), m.SegMeta.NextSeg)
		err := m.SyncMeta()
		if err != nil {
			log.Errorf("sync meta err:%v", err)
		}
		<-ticker.C
	}
}

// SyncMeta 同步元数据到文件
func (m *Meta) SyncMeta() error {
	err := writeMeta(m)
	if err != nil {
		return fmt.Errorf("writeSeg err: %v", err)
	}
	return nil
}
func readMeta(metaFile string) (*Meta, error) {
	metaByte, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, fmt.Errorf("read file err: %v", err)
	}
	h := new(Meta)
	err = json.Unmarshal(metaByte, &h)
	if err != nil {
		return nil, fmt.Errorf("ParseHeader err: %v", err)
	}
	log.Debugf("seg header :%v", h)
	h.path = metaFile
	return h, nil
}

func writeMeta(m *Meta) error {
	f, err := os.OpenFile(m.path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
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

// UpdateSegMeta --
func (m *Meta) UpdateSegMeta(segID segment.SegID, indexCount uint64) error {
	err := m.SegMeta.UpdateSegMeta(segID, indexCount)
	if err != nil {
		return fmt.Errorf("update seg meta err: %v", err)
	}
	err = m.SyncMeta()
	if err != nil {
		return fmt.Errorf("sync writeSeg err: %v", err)
	}
	return nil
}
