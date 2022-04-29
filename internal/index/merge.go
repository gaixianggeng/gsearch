package index

import (
	"doraemon/conf"
	"doraemon/internal/engine"
	"doraemon/internal/storage"
	"doraemon/pkg/utils"
	"encoding/json"
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
			err := m.merge(segs)
			if err != nil {
				log.Errorf("merge error: %v", err)
			}
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
func (m *MergeScheduler) merge(segs *MergeMessage) error {
	defer m.Done()

	log.Debugf("merge segs: %v", segs)

	// 获取merge的文件
	segList, docSize := m.getMergeFiles(segs)
	log.Debugf("prepare to merge seg list:%v", segList)

	// 恢复seg is_merging状态
	defer func() {
		for _, seg := range ([]*engine.SegInfo)(*segs) {
			// 如果merge失败，没有删除旧seg，需要恢复
			if s, ok := m.Meta.SegInfo[seg.SegID]; ok {
				s.IsMerging = false
			}
		}
	}()

	// 合并
	err := m.mergeSegments(segList, docSize)
	if err != nil {
		return err
	}
	return nil
}

// 合并k个升序链表 https://leetcode-cn.com/problems/merge-k-sorted-lists/
func (m *MergeScheduler) mergeSegments(segList []*segmentName, docSize uint64) error {
	// 初始化对应正排和倒排库
	segmentDBs := make([]segmentDB, 0)
	for _, seg := range segList {
		inDB := storage.NewInvertedDB(seg.term, seg.inverted)
		forDB := storage.NewForwardDB(seg.forward)
		segmentDBs = append(segmentDBs, segmentDB{inDB, forDB})
	}
	if len(segmentDBs) == 0 {
		log.Warn("no segment to merge")
		return nil
	}
	log.Debugf("final prepare to merge[%v]!", segList)

	termNodes := make([]*engine.TermNode, 0)
	for _, seg := range segmentDBs {
		termNode := new(engine.TermNode)
		termNode.Info = make(chan storage.TermInfo)
		termNode.DB = seg.inverted

		// 开启协程遍历读取
		go seg.inverted.GetTermCursor(termNode.Info)

		termNodes = append(termNodes, termNode)
	}

	// 合并
	res, err := engine.MergeKSegments(termNodes)
	if err != nil {
		log.Errorf("merge error: %v", err)
		return err
	}

	targetEng := engine.NewEngine(m.Meta, m.conf, engine.MergeMode)
	for token, pos := range res {
		c, _ := json.Marshal(pos)
		log.Infof("token:%s count:%d,pos:%s", token, pos.DocCount, c)
		err := targetEng.StoragePostings(pos)
		if err != nil {
			log.Errorf("storage postings err:%v", err)
			return err
		}
	}
	m.Meta.UpdateSegMeta(targetEng.CurrSegID, docSize)

	// 删除seginfo里的旧segemnts
	return nil
}

func (m *MergeScheduler) getMergeFiles(segs *MergeMessage) ([]*segmentName, uint64) {

	segList := make([]*segmentName, 0)
	docSize := uint64(0)
	for _, seg := range []*engine.SegInfo(*segs) {
		if seg.IsMerging {
			log.Infof("seg:%v is merging...", seg)
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
		docSize += seg.SegSize
	}
	return segList, docSize

}

// 判断seg是否存在，防止已经merge
func (m *MergeScheduler) segExists(termName, invertedName, forwardName string) bool {
	return utils.ExistFile(termName) && utils.ExistFile(invertedName) && utils.ExistFile(forwardName)
}

// NewScheduleer 创建调度器
func NewScheduleer(meta *engine.Meta, conf *conf.Config) *MergeScheduler {
	ch := make(chan *MergeMessage, conf.Merge.ChannelSize)

	// conf.Storage.Path = "../../data/"

	return &MergeScheduler{
		Message: ch,
		conf:    conf,
		Meta:    meta,
	}
}
