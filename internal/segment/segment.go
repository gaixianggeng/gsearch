package segment

import (
	"bytes"
	"doraemon/conf"
	"doraemon/internal/query"
	"doraemon/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Segment 段信息 封装term、倒排和正排库
type Segment struct {
	*storage.ForwardDB
	*storage.InvertedDB
	conf *conf.Config

	BufCount uint64 // 倒排索引缓冲区的文档数
	BufSize  uint64 // 设定的缓冲区大小

	PostingsHashBuf InvertedIndexHash // 倒排索引缓冲区
}

// Text2PostingsLists --
func (e *Segment) Text2PostingsLists(meta *SegMeta, tokens []query.Tokenization, docID uint64) error {

	bufInvertedHash := make(InvertedIndexHash)

	for _, token := range tokens {
		err := e.token2PostingsLists(bufInvertedHash, token.Token, token.Position, docID)
		if err != nil {
			return fmt.Errorf("text2PostingsLists token2PostingsLists err: %v", err)
		}
	}
	log.Debugf("bufInvertedHash:%v", bufInvertedHash)

	if e.PostingsHashBuf != nil && len(e.PostingsHashBuf) > 0 {
		// 合并命中相同的token的不同doc
		MergeInvertedIndex(e.PostingsHashBuf, bufInvertedHash)
	} else {
		e.PostingsHashBuf = make(InvertedIndexHash)
		e.PostingsHashBuf = bufInvertedHash
	}
	e.BufCount++
	return nil
}

// IsNeedFlush 是否需要落盘
func (e *Segment) IsNeedFlush() bool {
	if len(e.PostingsHashBuf) > 0 && (e.BufCount >= e.BufSize) {
		return true
	}
	return false
}

// token2PostingsLists --
func (e *Segment) token2PostingsLists(bufInvertHash InvertedIndexHash, token string,
	position uint64, docID uint64) error {

	// init
	bufInvert := new(InvertedIndexValue)

	if len(bufInvertHash) > 0 {
		if item, ok := bufInvertHash[token]; ok {
			bufInvert = item
		}
	}

	pl := new(PostingsList)
	if bufInvert != nil && bufInvert.PostingsList != nil {
		pl = bufInvert.PostingsList
		// 这里的positioinCount和下面bufInvert的positionCount是不一样的
		// 这里统计的是同一个docid的position的个数
		pl.PositionCount++
	} else {
		// 不为空表示写入操作，否则为查询
		termValue := new(storage.TermValue)
		if docID != 0 {
			termValue.DocCount = 1
			// docCount = 1
		} else {
			// docCount 用于召回排序使用
			var err error
			termValue, err = e.getTokenCount(token)
			if err != nil {
				return fmt.Errorf("token2PostingsLists GetTokenID err: %v", err)
			}

		}
		bufInvert = CreateNewInvertedIndex(token, termValue)
		bufInvertHash[token] = bufInvert
		pl = CreateNewPostingsList(docID)
		bufInvert.PostingsList = pl
	}
	// 存储位置信息
	pl.Positions = append(pl.Positions, position)
	// 统计该token关联的所有doc的position的个数
	bufInvert.PositionCount++

	return nil
}

// getTokenCount 通过token获取doc数量 insert 标识是写入还是查询 写入时不为空
func (e *Segment) getTokenCount(token string) (*storage.TermValue, error) {
	// _, c, err := e.FetchPostings(token)
	// if err != nil {
	// 	return 0, fmt.Errorf("getTokenCount FetchPostings err: %v", err)
	// }
	// return c, nil
	termInfo, err := e.InvertedDB.GetTermInfo(token)
	if err != nil || termInfo == nil {
		return nil, fmt.Errorf("getTokenCount GetTermInfo err: %v", err)
	}
	return termInfo, nil
}

// FetchPostings 通过token读取倒排表数据，返回倒排表、长度和err
func (e *Segment) FetchPostings(token string) (*PostingsList, uint64, error) {

	term, err := e.InvertedDB.GetTermInfo(token)
	if err != nil {
		return nil, 0, fmt.Errorf("FetchPostings getForwardAddr err: %v", err)
	}

	c, err := e.InvertedDB.GetInvertedDoc(term.Offset, term.Size)
	if err != nil {
		return nil, 0, fmt.Errorf("FetchPostings getDocInfo err: %v", err)
	}
	return decodePostings(bytes.NewBuffer(c))

}

// Flush 落盘操作
func (s *Segment) Flush() error {
	if len(s.PostingsHashBuf) == 0 {
		log.Warnf("Flush err: %v", "in.PostingsHashBuf is empty")
		return nil
	}
	log.Debugf("start storage...%v,len:%d", s.PostingsHashBuf, len(s.PostingsHashBuf))
	// title = ""表示文件读取结束
	for token, invertedIndex := range s.PostingsHashBuf {
		log.Debugf("token:%s,invertedIndex:%v\n", token, invertedIndex)
		err := s.storagePostings(invertedIndex)
		if err != nil {
			log.Errorf("updatePostings err: %v", err)
			return fmt.Errorf("updatePostings err: %v", err)
		}
	}
	return nil

}

// storagePostings 落盘
func (e *Segment) storagePostings(p *InvertedIndexValue) error {
	if p == nil {
		return fmt.Errorf("updatePostings p is nil")
	}

	// 编码
	buf, err := EncodePostings(p.PostingsList, p.DocCount)
	if err != nil {
		return fmt.Errorf("updatePostings encodePostings err: %v", err)
	}

	// 开始写入数据库
	return e.InvertedDB.StoragePostings(p.Token, buf.Bytes(), p.DocCount)
}

// Close --
func (e *Segment) Close() {
	e.InvertedDB.Close()
	e.ForwardDB.Close()
}

// NewSegments 创建新的segments 更新nextseg
func NewSegments(meta *SegMeta, conf *conf.Config, mode Mode) (SegID, map[SegID]*Segment) {

	segs := make(map[SegID]*Segment, 0)
	if mode == MergeMode || mode == IndexMode {
		segID := meta.NextSeg
		meta.NewSegmentItem()
		seg := NewSegment(segID, conf)
		segs[segID] = seg
		return segID, segs
	}

	for segID := range meta.SegInfo {
		seg := NewSegment(segID, conf)
		log.Infof("dbInit segID:%v,next:%v", segID, meta.NextSeg)
		segs[segID] = seg
	}
	return -1, segs

}

// NewSegment 创建新的segment
func NewSegment(segID SegID, conf *conf.Config) *Segment {

	postingsHashBuf := make(InvertedIndexHash)
	inDB, forDB := dbInit(segID, conf)
	return &Segment{
		PostingsHashBuf: postingsHashBuf,
		InvertedDB:      inDB,
		ForwardDB:       forDB,
		BufSize:         5,
		conf:            conf,
	}
}
