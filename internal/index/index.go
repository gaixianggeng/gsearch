package index

import (
	"brain/internal/storage"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// AddDocument 添加文档
func (e *Engine) AddDocument(doc *storage.Document) error {
	if doc.DocID > 0 && doc.Title != "" {
		err := e.forwardDB.Add(doc)
		if err != nil {
			return fmt.Errorf("AddDocument err: %v", err)
		}
		log.Debug(doc.DocID)
		err = e.text2PostingsLists(doc.DocID, []byte(doc.Title))
		if err != nil {
			return fmt.Errorf("text2postingslists err: %v", err)
		}
		e.bufCount++
		e.indexCount++
	}

	// 落盘操作 title = ""表示文件读取结束
	if len(e.postingsHashBuf) > 0 && (e.bufCount > e.bufSize || doc.Title == "") {

		for tokenID, invertedIndex := range e.postingsHashBuf {

			log.Debugf("tokenID:%d,invertedIndex:%v\n", tokenID, invertedIndex)
			e.updatePostings(invertedIndex)
		}

		// 重置
		e.postingsHashBuf = make(InvertedIndexHash)
		e.bufCount = 0
	}

	return nil

}

// 创建倒排列表
func createNewPostingsList(docID uint64) *PostingsList {
	p := new(PostingsList)
	p.DocID = docID
	p.positionCount = 1
	p.positions = make([]uint64, 0)
	return p
}

// 创建倒排索引
func createNewInvertedIndex(tokenID, docCount uint64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.docsCount = docCount
	p.TokenID = tokenID
	p.positionCount = 0
	p.postingsList = new(PostingsList)
	return p
}

// NewIndexEngine init
func NewIndexEngine() (*Engine, error) {
	return &Engine{
		bufSize: 30,
		N:       2,
	}, nil
}
