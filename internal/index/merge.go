package index

import (
	"doraemon/conf"
	"doraemon/internal/engine"
)

// MergeScheduler 合并调度器
type MergeScheduler struct {
}

// MergeQueue 合并队列
type MergeQueue []*engine.SegInfo

// 判断是否需要merge
func (m *MergeScheduler) mayMerge() {

}

// Merge 合并入口
func (m *MergeScheduler) merge() {

}

// Merge 合并segment
func mergeSegment() {

}

// NewScheduleer 创建调度器
func NewScheduleer(meta *engine.Meta, conf *conf.Config) *MergeScheduler {
	return &MergeScheduler{}
}
