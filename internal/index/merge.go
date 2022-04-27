package index

import (
	"doraemon/conf"
	"doraemon/internal/engine"
	"doraemon/internal/storage"
	"doraemon/pkg/utils"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// MergeScheduler 合并调度器
type MergeScheduler struct {
	Message chan *MergeMessage
	Meta    *engine.Meta
	conf    *conf.Config

	sync.WaitGroup
}

// 段初始化db
type segmentDB struct {
	inverted *storage.InvertedDB
	forward  *storage.ForwardDB
}

// 段名称
type segmentName struct {
	term     string
	inverted string
	forward  string
}

// MergeMessage 合并队列
type MergeMessage []*engine.SegInfo

// Merge 合并入口
func (m *MergeScheduler) Merge() {
	for {
		select {
		case segs := <-m.Message:
			log.Debugf("merge msg: %v", segs)
			// 合并
			m.merge(segs)

		default:
			log.Infof("sleep 1s...")
			time.Sleep(1e9)
		}

	}
}

//Close 关闭调度器
func (m *MergeScheduler) Close() {
	// 保证所有merge执行完毕
	m.Wait()

}

// 判断是否需要merge
// 通过meta数据中的seginfo来计算
func (m *MergeScheduler) mayMerge() {
	mess, isNeed := m.calculateSegs()
	if !isNeed {
		return
	}

	m.Add(1)
	m.Message <- mess
}

// 计算是否有段需要合并
func (m *MergeScheduler) calculateSegs() (*MergeMessage, bool) {
	segs := m.Meta.SegInfo
	log.Debugf("segs: %v", segs)

	// 判断是否需要合并

	segList := make([]*engine.SegInfo, 0)
	segList = append(segList, segs[0])
	segList = append(segList, segs[1])

	mes := MergeMessage(segList)
	return &mes, true
}

// Merge 合并segment
func (m *MergeScheduler) merge(segs *MergeMessage) {
	defer m.Done()

	log.Debugf("merge segs: %v", segs)

	// 获取merge的文件
	segList := m.getMergeFiles(segs)
	log.Debugf("prepare to merge seg list:%v", segList)

	// 初始化对应正排和倒排库
	segmentDBs := make([]*segmentDB, 0)
	for _, seg := range segList {
		inDB := storage.NewInvertedDB(seg.term, seg.inverted)
		forDB := storage.NewForwardDB(seg.forward)
		segmentDBs = append(segmentDBs, &segmentDB{inDB, forDB})
	}
	if len(segmentDBs) == 0 {
		log.Warn("no segment to merge")
		return
	}

	targetSeg := m.newSegment()

	// 合并
	m.mergeSegments(targetSeg, segmentDBs)
}

func (m *MergeScheduler) newSegment() *segmentDB {

	seg := m.Meta.NewSegment()
	log.Debugf("target seg id:%v, next id:%v", seg.SegID, m.Meta.NextSeg)

	path := m.conf.Storage.Path
	term := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.TermDBSuffix)
	inverted := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.InvertedDBSuffix)
	forward := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.ForwardDBSuffix)

	inDB := storage.NewInvertedDB(term, inverted)
	forDB := storage.NewForwardDB(forward)
	return &segmentDB{inDB, forDB}
}

// 合并k个升序链表 https://leetcode-cn.com/problems/merge-k-sorted-lists/
func (m *MergeScheduler) mergeSegments(targetDB *segmentDB, segmentDBs []*segmentDB) {
	log.Debugf("final prepare to merge!")

	for _, seg := range segmentDBs {
		seg.inverted.GetAllTerm()
	}
}

func (m *MergeScheduler) getMergeFiles(segs *MergeMessage) []*segmentName {

	segList := make([]*segmentName, 0)

	for _, seg := range []*engine.SegInfo(*segs) {
		if seg.IsMerging {
			continue
		}
		seg.IsMerging = true

		path := m.conf.Storage.Path
		term := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.TermDBSuffix)
		inverted := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.InvertedDBSuffix)
		forward := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.ForwardDBSuffix)

		if !m.segExists(term, inverted, forward) {
			continue
		}
		segName := new(segmentName)
		segName.forward = forward
		segName.inverted = inverted
		segName.term = term
		segList = append(segList, segName)

	}
	return segList

}

// 判断seg是否存在，防止已经merge
func (m *MergeScheduler) segExists(termName, invertedName, forwardName string) bool {
	return utils.ExistFile(termName) && utils.ExistFile(invertedName) && utils.ExistFile(forwardName)
}

// NewScheduleer 创建调度器
func NewScheduleer(meta *engine.Meta, conf *conf.Config) *MergeScheduler {
	ch := make(chan *MergeMessage, conf.Merge.ChannelSize)
	conf.Storage.Path = "../../data/"

	return &MergeScheduler{
		Message: ch,
		conf:    conf,
		Meta:    meta,
	}
}
