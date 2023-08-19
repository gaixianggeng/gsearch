package index

import (
	"fmt"
	"gsearch/conf"
	"gsearch/internal/engine"
	"gsearch/internal/segment"
	"gsearch/internal/storage"
)

// Index --
type Index struct {
	*engine.Engine
	*engine.Meta
	Conf *conf.Config
}

// AddDocument 添加文档
func (in *Index) AddDocument(doc *storage.Document) error {
	if doc == nil || doc.DocID <= 0 || doc.Title == "" {
		return fmt.Errorf("doc err: %v", "doc || doc_id || title is empty")
	}
	err := in.AddDoc(doc)
	if err != nil {
		return fmt.Errorf("forward doc add err: %v", err)
	}
	err = in.Text2PostingsLists(doc.Title, doc.DocID)
	if err != nil {
		return fmt.Errorf("text2postingslists err: %v", err)
	}
	return nil
}

// Close --
func (in *Index) Close() {
	in.Engine.Close()
}

// NewIndexEngine init
func NewIndexEngine(meta *engine.Meta, c *conf.Config) (*Index, error) {
	e := engine.NewEngine(meta, c, segment.IndexMode)
	return &Index{
		Engine: e,
		Conf:   c,
		Meta:   meta,
	}, nil
}
