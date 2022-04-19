package recall

import (
	"brain/internal/engine"

	log "github.com/sirupsen/logrus"
)

// Recall 查询召回
type Recall struct {
	*engine.Engine
	queryTokenHash *engine.InvertedIndexHash
}

// SearchResult 查询结果
type SearchResult struct {
}

// Search 入口
func (r *Recall) Search(query string) *SearchResult {
	return nil
}

func (r *Recall) splitQuery2Tokens(query string) {
	err := r.Text2PostingsLists(query, 0)
	if err != nil {
		log.Errorf("text2postingslists err: %v", err)
		return
	}
	r.queryTokenHash = new(engine.InvertedIndexHash)
	*r.queryTokenHash = r.Engine.PostingsHashBuf
	log.Debugf("queryHash:%v,engine:%v", r.queryTokenHash, &r.Engine.PostingsHashBuf)

}

// NewRecall new
func NewRecall(e *engine.Engine) *Recall {
	return &Recall{e, nil}
}
