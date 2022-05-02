package engine

import (
	"doraemon/conf"
	"doraemon/internal/query"
	"doraemon/internal/segment"
	"doraemon/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Engine 写入引擎
type Engine struct {
	CurrSegID segment.SegID //当前engine关联的segID 查询时为-1
	Seg       map[segment.SegID]*segment.Segment
	meta      *Meta // 元数据
	conf      *conf.Config
	scheduler *MergeScheduler

	N int32 // ngram

}

// AddDoc 添加正排
func (e *Engine) AddDoc(doc *storage.Document) error {
	return e.Seg[e.CurrSegID].AddForward(doc)
}

// Text2PostingsLists --
func (e *Engine) Text2PostingsLists(text string, docID uint64) error {

	go e.scheduler.Merge()
	tokens, err := query.Ngram(text, e.N)
	if err != nil {
		return fmt.Errorf("text2PostingsLists Ngram err: %v", err)
	}
	err = e.Seg[e.CurrSegID].Text2PostingsLists(e.meta.SegMeta, tokens, docID)
	if err != nil {
		return fmt.Errorf("text2PostingsLists err: %v", err)
	}
	// 达到阈值
	if e.Seg[e.CurrSegID].IsNeedFlush() {
		log.Infof("text2PostingsLists need flush")
		e.Flush()
	}
	e.indexCount()
	return nil
}

// Flush --
func (e *Engine) Flush(isEnd ...bool) error {

	e.Seg[e.CurrSegID].Flush()

	// update meta info
	err := e.meta.UpdateSegMeta(e.CurrSegID, e.Seg[e.CurrSegID].BufCount)
	if err != nil {
		log.Errorf("update seg meta err:%v", err)
		return err
	}

	e.Seg[e.CurrSegID].Close()
	delete(e.Seg, e.CurrSegID)

	if len(e.meta.SegMeta.SegInfo) > 1 {
		e.scheduler.mayMerge()
	}

	// new
	if len(isEnd) > 0 && isEnd[0] {
		return nil
	}
	segID, seg := segment.NewSegments(e.meta.SegMeta, e.conf, segment.IndexMode)
	e.CurrSegID = segID
	e.Seg = seg
	return nil

}

// UpdateCount --
func (e *Engine) UpdateCount(num uint64) error {
	seg := e.Seg[e.CurrSegID]
	count, err := seg.ForwardCount()
	if err != nil {
		if err.Error() == ErrCountKeyNotFound {
			count = 0
		} else {
			return fmt.Errorf("updateCount err: %v", err)
		}
	}
	count += num
	return seg.UpdateForwardCount(count)
}

func (e *Engine) indexCount() {
	e.meta.Lock()
	e.meta.IndexCount++
	e.meta.Unlock()
}

// func (e *Engine) Flush() error {
// 	return e.Seg[e.CurrSegID].Flush()
// }

// StoragePostings 落盘
func (e *Engine) StoragePostings(p *segment.InvertedIndexValue) error {
	if p == nil {
		return fmt.Errorf("updatePostings p is nil")
	}

	// 编码
	buf, err := segment.EncodePostings(p.PostingsList, p.DocCount)
	if err != nil {
		return fmt.Errorf("updatePostings encodePostings err: %v", err)
	}

	// 开始写入数据库
	return e.Seg[e.CurrSegID].InvertedDB.StoragePostings(p.Token, buf.Bytes(), p.DocCount)
}

// Close --
func (e *Engine) Close() {
	for _, seg := range e.Seg {
		seg.Close()
	}

	e.scheduler.Close()
}

// NewEngine --
// 每次初始化的时候调整meta数据
func NewEngine(meta *Meta, conf *conf.Config, engineMode segment.Mode) *Engine {

	sche := NewScheduleer(meta, conf)
	segID, seg := segment.NewSegments(meta.SegMeta, conf, engineMode)
	return &Engine{
		CurrSegID: segID,
		Seg:       seg,
		N:         2,
		meta:      meta,
		conf:      conf,
		scheduler: sche,
	}
}
