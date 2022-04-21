package recall

import (
	"brain/internal/engine"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

// Recall 查询召回
type Recall struct {
	*engine.Engine
	queryToken []*queryTokenHash
}

// 用于实现排序map
type queryTokenHash struct {
	token         string
	invertedIndex *engine.InvertedIndexValue
}

// SearchResult 查询结果
type SearchResult struct {
}

// 游标 标识当前位置
type searchCursor struct {
	doc     *engine.PostingsList
	current *engine.PostingsList
}

// Search 入口
func (r *Recall) Search(query string) (*SearchResult, error) {
	err := r.splitQuery2Tokens(query)
	if err != nil {
		log.Errorf("splitQuery2Tokens err: %v", err)
		return nil, fmt.Errorf("splitQuery2Tokens err: %v", err)
	}
	return r.searchDoc()
}

func (r *Recall) splitQuery2Tokens(query string) error {
	err := r.Text2PostingsLists(query, 0)
	if err != nil {
		return fmt.Errorf("text2postingslists err: %v", err)
	}
	log.Debugf("queryHash:%v,engine:%v", r.queryToken, &r.Engine.PostingsHashBuf)
	return nil
}

func (r *Recall) searchDoc() (*SearchResult, error) {

	r.sortToken(r.Engine.PostingsHashBuf)
	if len(r.queryToken) == 0 {
		return nil, fmt.Errorf("queryTokenHash is nil")
	}

	return nil, nil
}

// token 根据doc count升序排序
func (r *Recall) sortToken(postHash engine.InvertedIndexHash) {
	tokenHash := make([]*queryTokenHash, 0)
	for token, invertedIndex := range postHash {
		q := new(queryTokenHash)
		q.token = token
		q.invertedIndex = invertedIndex
		tokenHash = append(tokenHash, q)
	}

	log.Debugf("tokenHash:%v", tokenHash)
	sort.Sort(docCountSort(tokenHash))
	log.Debugf("tokenHash:%v", tokenHash)
	r.queryToken = tokenHash
	return
}

// NewRecall new
func NewRecall(e *engine.Engine) *Recall {
	return &Recall{e, nil}
}
