package index

import (
	"fmt"
	"gsearch/conf"
	"gsearch/internal/engine"
	"gsearch/internal/meta"
	"gsearch/internal/segment"
	"gsearch/internal/storage"
	"gsearch/pkg/utils/log"
	"strconv"
	"strings"
)

// Index --
type Index struct {
	engine *engine.Engine
	*meta.Profile
	Conf *conf.Config
}

func (in *Index) IndexDoc() {
	log.Infof("addDoc start")
	docList := readFiles(in.Conf.Source.Files)
	go in.engine.Scheduler.Merge()
	for _, item := range docList {
		doc, err := in.doc2Struct(item)
		if err != nil {
			log.Errorf("doc2Struct err: %v", err)
			doc = new(storage.Document)
		}
		log.Debugf("doc_id:%v,title:%s", doc.DocID, doc.Title)
		err = in.addDoc(doc)
		if err != nil {
			log.Errorf("AddDocument err: %v", err)
			break
		}
	}
	// 读取结束 写入磁盘
	in.engine.Flush(true)
}

// addDoc 添加文档
func (in *Index) addDoc(doc *storage.Document) error {
	if doc == nil || doc.DocID <= 0 || doc.Title == "" {
		return fmt.Errorf("doc err: %v", "doc || doc_id || title is empty")
	}
	// 添加正排
	err := in.engine.AddDoc(doc)
	if err != nil {
		return fmt.Errorf("forward doc add err: %v", err)
	}
	err = in.engine.Text2PostingsLists(doc.Title, doc.DocID)
	if err != nil {
		return fmt.Errorf("text2postingslists err: %v", err)
	}
	return nil
}
func (in *Index) doc2Struct(docStr string) (*storage.Document, error) {
	d := strings.Split(docStr, "\t")
	if len(d) < 3 {
		return nil, fmt.Errorf("doc2Struct err: %v", "docStr is not right")
	}
	doc := new(storage.Document)

	docID, err := strconv.Atoi(d[0])
	if err != nil {
		return nil, err
	}
	doc.DocID = uint64(docID)
	doc.Title = d[1]
	doc.Body = d[2]
	return doc, nil
}

// Close --
func (in *Index) Close() {
	in.engine.Close()
}

// NewIndexEngine init
func NewIndexEngine(meta *meta.Profile, c *conf.Config) (*Index, error) {
	e := engine.NewEngine(meta, c, segment.IndexMode)
	return &Index{
		engine:  e,
		Conf:    c,
		Profile: meta,
	}, nil
}
