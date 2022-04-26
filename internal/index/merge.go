package index

import (
	"doraemon/conf"
	"doraemon/internal/engine"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// MergeScheduler 合并调度器
type MergeScheduler struct {
	Message chan *MergeMessage
	sync.WaitGroup
}

// MergeMessage 合并队列
type MergeMessage []*engine.SegInfo

// 判断是否需要merge
func (m *MergeScheduler) mayMerge(segInfo *engine.SegInfo) {

}

// Merge 合并入口
func (m *MergeScheduler) merge() {
	select {
	case msg := <-m.Message:
		log.Debugf("merge msg: %v", msg)
		// 合并
	default:
		time.Sleep(1e9)
	}

}

// Merge 合并segment
func mergeSegment() {

}

// NewScheduleer 创建调度器
func NewScheduleer(meta *engine.Meta, conf *conf.Config) *MergeScheduler {
	ch := make(chan *MergeMessage, conf.Merge.ChannelSize)
	return &MergeScheduler{
		Message: ch,
	}
}
