package recall

import "brain/internal/engine"

// Recall 查询召回
type Recall struct {
	*engine.Engine
}

// SearchResult 查询结果
type SearchResult struct {
}

// Search 入口
func (r *Recall) Search(query string) *SearchResult {
	return nil
}

func (r *Recall) splitQuery2Tokens(query string) {
	r.Text2PostingsLists(query, 0)
}

// NewRecall new
func NewRecall(e *engine.Engine) *Recall {
	return &Recall{e}
}
