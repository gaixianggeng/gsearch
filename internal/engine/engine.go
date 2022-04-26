package engine

import (
	"bytes"
	"doraemon/conf"
	"doraemon/internal/query"
	"doraemon/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	termDBSuffix     = ".term"
	invertedDBSuffix = ".inverted"
	forwardDBSuffix  = ".forward"
)
var (
	termName     = ""
	invertedName = ""
	forwardName  = ""
)

// Engine 写入引擎
type Engine struct {
	ForwardFileName string

	ForwardDB  *storage.ForwardDB
	InvertedDB *storage.InvertedDB

	Meta *Meta

	PostingsHashBuf InvertedIndexHash // 倒排索引缓冲区
	BufCount        uint64            //倒排索引缓冲区的文档数
	BufSize         uint64
	N               int32 // ngram
}

// PostingsList 倒排列表
type PostingsList struct {
	DocID         uint64
	Positions     []uint64
	PositionCount uint64
	Next          *PostingsList
}

//InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string
	PostingsList  *PostingsList
	DocCount      uint64
	PositionCount uint64 // 查询使用，写入的时候暂时用不到
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*InvertedIndexValue

// CreateNewPostingsList 创建倒排列表
func CreateNewPostingsList(docID uint64) *PostingsList {
	p := new(PostingsList)
	p.DocID = docID
	p.PositionCount = 1
	p.Positions = make([]uint64, 0)
	return p
}

// CreateNewInvertedIndex 创建倒排索引
func CreateNewInvertedIndex(token string, docCount uint64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.DocCount = docCount
	p.Token = token
	p.PositionCount = 0
	p.PostingsList = new(PostingsList)
	return p
}

// Close --
func (e *Engine) Close() {
	e.InvertedDB.Close()
	e.ForwardDB.Close()
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
		return nil, 0, fmt.Errorf("FetchPostings getForwordAddr err: %v", err)
	}

	c, err := e.InvertedDB.GetForwordContent(term[1], term[2])
	if err != nil {
		return nil, 0, fmt.Errorf("FetchPostings getForwordContent err: %v", err)
	}
	return decodePostings(bytes.NewBuffer(c))
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
func NewEngine(meta *Meta, conf *conf.Config) *Engine {
	dbInit(meta, conf)
	inverted := storage.NewInvertedDB(
		termName, invertedName)
	forward := storage.NewForwardDB(forwardName)

	return &Engine{
		ForwardFileName: forwardName,
		Meta:            meta,
		InvertedDB:      inverted,
		ForwardDB:       forward,
		BufSize:         1000,
		N:               2,
	}

}
func dbInit(meta *Meta, conf *conf.Config) error {
	// 获取最新的segment id
	newSeg := meta.NextSeg
	termName = fmt.Sprintf("%s%d%s", conf.Storage.Path, newSeg, termDBSuffix)
	invertedName = fmt.Sprintf("%s%d%s", conf.Storage.Path, newSeg, invertedDBSuffix)
	forwardName = fmt.Sprintf("%s%d%s", conf.Storage.Path, newSeg, forwardDBSuffix)
	meta.NextSeg++
	return nil
}
