package engine

import (
	"doraemon/conf"
	"doraemon/internal/segment"
	"doraemon/internal/storage"
	"doraemon/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// MergeScheduler 合并调度器
type MergeScheduler struct {
	Message chan *MergeMessage
	Meta    *Meta
	conf    *conf.Config

	sync.WaitGroup
}

// // 段初始化db
// type segmentDB struct {
// 	inverted *storage.InvertedDB
// 	forward  *storage.ForwardDB
// }

// 段名称
// type segmentName struct {
// 	term     string
// 	inverted string
// 	forward  string
// }

// MergeMessage 合并队列
type MergeMessage []*segment.SegInfo

// Merge 合并入口
func (m *MergeScheduler) Merge() {
	for {
		select {
		case segs := <-m.Message:
			log.Infof("Merge msg: %v", segs)
			// 合并
			err := m.merge(segs)
			if err != nil {
				log.Errorf("merge error: %v", err)
			}
		case <-time.After(1e9):
			log.Infof("sleep 1s...")
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
	// 已存在超过2个segment，则需要判断seg是否需要merge
	if len(m.Meta.SegMeta.SegInfo) <= 1 {
		log.Infof("seg count: %v, no need merge", len(m.Meta.SegMeta.SegInfo))
		return
	}

	mess, isNeed := m.calculateSegs()
	if !isNeed {
		return
	}

	m.Add(1)
	m.Message <- mess

	log.Infof("merge segs: %v", mess)
}

// 计算是否有段需要合并
func (m *MergeScheduler) calculateSegs() (*MergeMessage, bool) {
	segs := m.Meta.SegMeta.SegInfo
	log.Debugf("segs: %v", segs)

	// 判断是否需要合并

	segList := make([]*segment.SegInfo, 0)
	segList = append(segList, segs[0])
	segList = append(segList, segs[1])

	mes := MergeMessage(segList)
	return &mes, true
}

// Merge 合并segment
func (m *MergeScheduler) merge(segs *MergeMessage) error {
	defer m.Done()

	log.Debugf("merge segs: %v", segs)

	// 恢复seg is_merging状态
	defer func() {
		for _, seg := range ([]*segment.SegInfo)(*segs) {
			// 如果merge失败，没有删除旧seg，需要恢复
			if s, ok := m.Meta.SegMeta.SegInfo[seg.SegID]; ok {
				s.IsMerging = false
			}
		}
	}()

	// 合并
	err := m.mergeSegments(segs)
	if err != nil {
		log.Errorf("merge error: %v", err)
		return err
	}
	return nil
}

// 合并k个升序链表 https://leetcode-cn.com/problems/merge-k-sorted-lists/
// term表需要合并k个升序，以及处理对应的倒排数据
// 正排表直接merge即可
func (m *MergeScheduler) mergeSegments(segs *MergeMessage) error {
	// 获取merge的文件
	// segMap, docSize := m.getMergeFiles(segs)
	// log.Debugf("prepare to merge seg list:%v,docsize:%d", segMap, docSize)

	// 初始化对应正排和倒排库
	segmentDBs := make([]*segment.Segment, 0)
	docSize := uint64(0)
	for _, segInfo := range []*segment.SegInfo(*segs) {
		docSize += segInfo.SegSize
		s := segment.NewSegment(segInfo.SegID, m.conf)
		segmentDBs = append(segmentDBs, s)
	}
	if len(segmentDBs) == 0 {
		log.Warn("no segment to merge")
		return nil
	}

	termNodes := make([]*segment.TermNode, 0)
	termChs := make([]chan storage.KvInfo, 0)

	forNodes := make([]*segment.TermNode, 0)
	forChs := make([]chan storage.KvInfo, 0)
	for _, seg := range segmentDBs {
		termNode := new(segment.TermNode)
		termNode.DB = seg

		// 开启协程遍历读取
		termCh := make(chan storage.KvInfo)
		go seg.GetInvertedTermCursor(termCh)

		forCh := make(chan storage.KvInfo)
		go seg.GetForwardCursor(forCh)

		termNodes = append(termNodes, termNode)
		termChs = append(termChs, termCh)

		forNodes = append(forNodes, new(segment.TermNode))
		forChs = append(forChs, forCh)
	}

	// 合并term和倒排数据,返回合并后的数据
	res, err := segment.MergeKTermSegments(termNodes, termChs)
	if err != nil {
		log.Errorf("merge error: %v", err)
		return err
	}

	targetEng := NewEngine(m.Meta, m.conf, segment.MergeMode)

	// 落盘
	for token, pos := range res {
		c, _ := json.Marshal(pos)
		log.Infof("token:%s count:%d,pos:%s", token, pos.DocCount, c)
		err := targetEng.StoragePostings(pos)
		if err != nil {
			log.Errorf("storage postings err:%v", err)
			return err
		}
	}

	log.Debugf("start forwatd:%s", strings.Repeat("-", 20))

	// 合并正排数据
	err = segment.MergeKForwardSegments(targetEng.Seg[targetEng.CurrSegID], forNodes, forChs)
	if err != nil {
		log.Errorf("forward merge error: %v", err)
		return err
	}

	// update meta info
	err = m.Meta.UpdateSegMeta(targetEng.CurrSegID, docSize)
	if err != nil {
		log.Errorf("update seg meta err:%v", err)
		return err
	}

	// delete old segs
	err = m.deleteOldSeg([]*segment.SegInfo(*segs))
	if err != nil {
		log.Errorf("delete old seg error: %v", err)
		return err
	}
	return nil
}

func (m *MergeScheduler) deleteOldSeg(segInfos []*segment.SegInfo) error {

	for _, segInfo := range segInfos {
		if s, ok := m.Meta.SegMeta.SegInfo[segInfo.SegID]; ok {
			s.IsMerging = false
			delete(m.Meta.SegMeta.SegInfo, segInfo.SegID)
			err := m.deleteSegFile(segInfo.SegID)
			if err != nil {
				log.Errorf("delete old seg error: %v", err)
				return err
			}
		} else {
			return fmt.Errorf("delete old seg error: %v", segInfo)
		}
	}
	return nil
}

func (m *MergeScheduler) deleteSegFile(segID segment.SegID) error {
	term, inverted, forward := segment.GetDBName(m.conf, segID)

	log.Debugf("delete seg file forward:%s,invert:%s,term:%s",
		forward, inverted, term)
	err := os.Remove(inverted)
	if err != nil {
		return err
	}
	os.Remove(term)
	if err != nil {
		return err
	}
	os.Remove(forward)
	if err != nil {
		return err
	}
	return nil

}

// func (m *MergeScheduler) getMergeFiles(segs *MergeMessage) (map[segment.SegID]*segmentName, uint64) {

// 	segMap := make(map[segment.SegID]*segmentName, 0)
// 	docSize := uint64(0)
// 	for _, seg := range []*segment.SegInfo(*segs) {
// 		if seg.IsMerging {
// 			log.Infof("seg:%v is merging...", seg)
// 			continue
// 		}
// 		seg.IsMerging = true

// 		path := m.conf.Storage.Path
// 		term := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.TermDBSuffix)
// 		inverted := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.InvertedDBSuffix)
// 		forward := fmt.Sprintf("%s%d%s", path, seg.SegID, engine.ForwardDBSuffix)

// 		if !m.segExists(term, inverted, forward) {
// 			continue
// 		}
// 		segName := new(segmentName)
// 		segName.forward = forward
// 		segName.inverted = inverted
// 		segName.term = term
// 		segMap[seg.SegID] = segName

// 		docSize += seg.SegSize
// 	}
// 	return segMap, docSize

// }

// 判断seg是否存在，防止已经merge
func (m *MergeScheduler) segExists(termName, invertedName, forwardName string) bool {
	return utils.ExistFile(termName) && utils.ExistFile(invertedName) && utils.ExistFile(forwardName)
}

// NewScheduleer 创建调度器
func NewScheduleer(meta *Meta, conf *conf.Config) *MergeScheduler {
	ch := make(chan *MergeMessage, conf.Merge.ChannelSize)

	// conf.Storage.Path = "../../data/"

	return &MergeScheduler{
		Message: ch,
		conf:    conf,
		Meta:    meta,
	}
}
