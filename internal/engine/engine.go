package engine

import (
	"bytes"
	"doraemon/conf"
	"doraemon/internal/query"
	"doraemon/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Engine 写入引擎
type Engine struct {
	Meta *Meta

	// MaxSegmentCount int64 // 最大segment数,超出就要merge

	ForwardDB  *storage.ForwardDB
	InvertedDB *storage.InvertedDB

	PostingsHashBuf InvertedIndexHash // 倒排索引缓冲区
	BufCount        uint64            //倒排索引缓冲区的文档数
	BufSize         uint64

	// query
	N int32 // ngram
}

// Text2PostingsLists --
func (e *Engine) Text2PostingsLists(text string, docID uint64) error {
	tokens, err := query.Ngram(text, e.N)
	if err != nil {
		return fmt.Errorf("text2PostingsLists Ngram err: %v", err)
	}
	bufInvertedHash := make(InvertedIndexHash)

	for _, token := range tokens {
		err := e.Token2PostingsLists(bufInvertedHash, token.Token, token.Position, docID)
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
	return nil
}

// Token2PostingsLists --
func (e *Engine) Token2PostingsLists(bufInvertHash InvertedIndexHash, token string,
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
		docCount := uint64(0)
		if docID != 0 {
			docCount = 1
		} else {
			// docCount 用于召回排序使用
			var err error
			docCount, err = e.getTokenCount(token)
			if err != nil {
				return fmt.Errorf("token2PostingsLists GetTokenID err: %v", err)
			}

		}
		bufInvert = CreateNewInvertedIndex(token, docCount)
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

// FetchPostings 通过token读取倒排表数据，返回倒排表、长度和err
func (e *Engine) FetchPostings(token string) (*PostingsList, uint64, error) {

	term, err := e.InvertedDB.GetTermInfo(token)
	if err != nil {
		return nil, 0, fmt.Errorf("FetchPostings getForwardAddr err: %v", err)
	}

	c, err := e.InvertedDB.GetDocInfo(term[1], term[2])
	if err != nil {
		return nil, 0, fmt.Errorf("FetchPostings getDocInfo err: %v", err)
	}
	return decodePostings(bytes.NewBuffer(c))
}

// Close --
func (e *Engine) Close() {
	e.InvertedDB.Close()
	e.ForwardDB.Close()
}

// getTokenCount 通过token获取doc数量 insert 标识是写入还是查询 写入时不为空
func (e *Engine) getTokenCount(token string) (uint64, error) {
	// _, c, err := e.FetchPostings(token)
	// if err != nil {
	// 	return 0, fmt.Errorf("getTokenCount FetchPostings err: %v", err)
	// }
	// return c, nil
	termInfo, err := e.InvertedDB.GetTermInfo(token)
	if err != nil || termInfo == nil {
		return 0, fmt.Errorf("getTokenCount GetTermInfo err: %v", err)
	}
	return termInfo[0], nil
}

// NewEngine --
func NewEngine(meta *Meta, conf *conf.Config, engineMode Mode) *Engine {
	inDB, forDB := dbInit(meta, conf, engineMode)
	return &Engine{
		Meta:       meta,
		InvertedDB: inDB,
		ForwardDB:  forDB,
		BufSize:    5,
		N:          2,
	}

}

// 读取对应的segment文件下的db
func dbInit(meta *Meta, conf *conf.Config, mode Mode) (*storage.InvertedDB, *storage.ForwardDB) {
	segID := uint64(0)
	if mode == SearchMode {
		for _, seg := range meta.SegInfo {
			// 检查是否可读
			if !seg.IsReading {
				segID = seg.SegID
				seg.IsReading = true
				break
			}
		}
	} else if mode == IndexMode {
		segID = meta.NextSeg
	} else {
		log.Fatalf("dbInit mode err: %v", mode)
	}
	termName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, TermDBSuffix)
	invertedName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, InvertedDBSuffix)
	forwardName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, ForwardDBSuffix)
	log.Debugf(
		"index:[termName:%s,invertedName:%s,forwardName:%s]",
		termName,
		invertedName,
		forwardName,
	)
	return storage.NewInvertedDB(termName, invertedName), storage.NewForwardDB(forwardName)
}
