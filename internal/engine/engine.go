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
	meta      *Meta // 元数据
	conf      *conf.Config
	Scheduler *MergeScheduler

	BufCount        uint64                             // 倒排索引缓冲区的文档数
	BufSize         uint64                             // 设定的缓冲区大小
	PostingsHashBuf segment.InvertedIndexHash          // 倒排索引缓冲区
	CurrSegID       segment.SegID                      // 当前engine关联的segID 查询时为-1
	Seg             map[segment.SegID]*segment.Segment // 当前engine关联的segment

	N int32 // ngram

}

// AddDoc 添加正排
func (e *Engine) AddDoc(doc *storage.Document) error {
	return e.Seg[e.CurrSegID].AddForward(doc)
}

// Text2PostingsLists 文本转倒排索引
func (e *Engine) Text2PostingsLists(text string, docID uint64) error {
	tokens, err := query.NGram(text, e.N)
	if err != nil {
		return fmt.Errorf("text2PostingsLists Ngram err: %v", err)
	}

	bufInvertedHash := make(segment.InvertedIndexHash)

	for _, token := range tokens {
		err := segment.Token2PostingsLists(bufInvertedHash, token.Token, token.Position, docID)
		if err != nil {
			return fmt.Errorf("text2PostingsLists token2PostingsLists err: %v", err)
		}
	}
	log.Debugf("bufInvertedHash:%v", bufInvertedHash)

	if e.PostingsHashBuf != nil && len(e.PostingsHashBuf) > 0 {
		// 合并命中相同的token的不同doc
		segment.MergeInvertedIndex(e.PostingsHashBuf, bufInvertedHash)
	} else {
		// 已经初始化过了
		e.PostingsHashBuf = bufInvertedHash
	}

	e.BufCount++

	// 达到阈值
	if len(e.PostingsHashBuf) > 0 && (e.BufCount >= e.BufSize) {
		log.Infof("text2PostingsLists need flush")
		e.Flush()
	}

	e.indexCount()
	return nil
}

// Flush --
func (e *Engine) Flush(isEnd ...bool) error {

	e.Seg[e.CurrSegID].Flush(e.PostingsHashBuf)

	// update meta info
	err := e.meta.UpdateSegMeta(e.CurrSegID, e.BufCount)
	if err != nil {
		log.Errorf("update seg meta err:%v", err)
		return err
	}

	e.UpdateCount(e.meta.IndexCount)
	e.Seg[e.CurrSegID].Close()
	delete(e.Seg, e.CurrSegID)

	if len(e.meta.SegMeta.SegInfo) > 1 {
		e.Scheduler.MayMerge()
	}

	// new
	if len(isEnd) > 0 && isEnd[0] {
		return nil
	}
	segID, seg := segment.NewSegments(e.meta.SegMeta, e.conf, segment.IndexMode)

	e.BufCount = 0
	e.PostingsHashBuf = make(segment.InvertedIndexHash)
	e.CurrSegID = segID
	e.Seg = seg
	return nil

}

// UpdateCount 更新文档数量
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

// indexCount index计数
func (e *Engine) indexCount() {
	e.meta.Lock()
	e.meta.IndexCount++
	e.meta.Unlock()
}

// // StoragePostings 落盘
// func (e *Engine) StoragePostings(p *segment.InvertedIndexValue) error {
// 	if p == nil {
// 		return fmt.Errorf("updatePostings p is nil")
// 	}

// 	// 编码
// 	buf, err := segment.EncodePostings(p.PostingsList, p.DocCount)
// 	if err != nil {
// 		return fmt.Errorf("updatePostings encodePostings err: %v", err)
// 	}

// 	// 开始写入数据库
// 	return e.Seg[e.CurrSegID].InvertedDB.StoragePostings(p.Token, buf.Bytes(), p.DocCount)
// }

// Close --
func (e *Engine) Close() {
	for _, seg := range e.Seg {
		seg.Close()
	}

	e.Scheduler.Close()
}

// NewEngine --
// 每次初始化的时候调整meta数据
func NewEngine(meta *Meta, conf *conf.Config, engineMode segment.Mode) *Engine {

	sche := NewScheduleer(meta, conf)
	segID, seg := segment.NewSegments(meta.SegMeta, conf, engineMode)
	return &Engine{
		CurrSegID:       segID,
		Seg:             seg,
		N:               2,
		meta:            meta,
		conf:            conf,
		Scheduler:       sche,
		PostingsHashBuf: make(segment.InvertedIndexHash),
		BufSize:         5,
	}
}
