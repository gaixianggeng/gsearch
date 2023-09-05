package engine

import "gsearch/internal/document"

// Engine represents the search engine core
type Engine struct{}

// Index indexes a document
func (e *Engine) Index(doc document.Document) {
	// TODO: Implement the indexing logic
}

// Search performs a search based on the given query
func (e *Engine) Search(query string) ([]interface{}, error) {
	// TODO: Implement the search logic based on the query
	return nil, nil
}

func New() *Engine {
	return &Engine{}
}
