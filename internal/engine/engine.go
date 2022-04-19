package engine

import (
	"brain/internal/query"
	"brain/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Engine 写入引擎
type Engine struct {
	ForwardDB  *storage.ForwardDB
	InvertedDB *storage.InvertedDB

	PostingsHashBuf InvertedIndexHash // 倒排索引缓冲区
	BufCount        uint64            //倒排索引缓冲区的文档数
	BufSize         uint64
	IndexCount      uint64
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
	DocsCount     uint64
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
	p.DocsCount = docCount
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
		log.Debug("mergeInvertedIndex-----")
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

	// doc_id用来标识写入数据还是查询数据
	// ?? docCount 应该是用于查询
	docCount, err := e.InvertedDB.GetTokenCount(token, docID)
	if err != nil {
		return fmt.Errorf("token2PostingsLists GetTokenID err: %v", err)
	}

	if len(bufInvertHash) > 0 {
		if item, ok := bufInvertHash[token]; ok {
			bufInvert = item
		}
	}

	pl := new(PostingsList)
	if bufInvert != nil && bufInvert.PostingsList != nil {
		log.Debug("bufInvert.postingsList is not nil")
		pl = bufInvert.PostingsList
		// 这里的positioinCount和下面bufInvert的positionCount是不一样的
		// 这里统计的是同一个docid的position的个数
		pl.PositionCount++
	} else {
		log.Debug("bufInvert.postingsList is nil")
		if docID != 0 {
			docCount = 1
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

// NewEngine --
func NewEngine(termDB, invertedDB, forwardDB string) *Engine {
	inverted := storage.NewInvertedDB(
		termDB, invertedDB)
	forward := storage.NewForwardDB(forwardDB)

	return &Engine{
		InvertedDB: inverted,
		ForwardDB:  forward,
		BufSize:    1,
		N:          2,
	}

}
