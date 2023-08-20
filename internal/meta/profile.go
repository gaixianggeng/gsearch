package meta

import (
	"encoding/json"
	"fmt"
	"gsearch/internal/segment"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Profile 元数据
type Profile struct {
	sync.RWMutex
	Version    string           `json:"version"` // 版本号
	IndexCount uint64           `json:"index"`
	SegMeta    *segment.SegMeta `json:"seg_meta"`
	path       string           `json:"-"` // 元数据文件路径
}

// SyncByTicker 定时同步元数据
func (m *Profile) SyncByTicker(ticker *time.Ticker) {
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
func (m *Profile) SyncMeta() error {
	err := m.writeMeta()
	if err != nil {
		return fmt.Errorf("writeSeg err: %v", err)
	}
	return nil
}

func (m *Profile) writeMeta() error {
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
func (m *Profile) UpdateSegMeta(segID segment.SegID, indexCount uint64) error {
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
