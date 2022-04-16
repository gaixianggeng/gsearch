package index

import (
	"brain/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// AddDocument 添加文档
func (e *Engine) AddDocument(doc *storage.Document) error {
	if doc.DocID > 0 && doc.Title != "" {
		// err := e.forwardDB.Add(doc)

		err := e.text2PostingsLists(doc.DocID, (doc.Title))
		if err != nil {
			return fmt.Errorf("text2postingslists err: %v", err)
		}
		e.bufCount++
		e.indexCount++
	}

	log.Debugf("start storage...%v,len:%d", e.postingsHashBuf, len(e.postingsHashBuf))

	// 落盘操作 title = ""表示文件读取结束
	if len(e.postingsHashBuf) > 0 && (e.bufCount > e.bufSize || doc.Title == "") {

		for token, invertedIndex := range e.postingsHashBuf {

			log.Debugf("token:%s,invertedIndex:%v\n", token, invertedIndex)
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
func createNewInvertedIndex(token string, docCount uint64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.docsCount = docCount
	p.Token = token
	p.positionCount = 0
	p.postingsList = new(PostingsList)
	return p
}

// NewIndexEngine init
func NewIndexEngine(termDB, forwardDB string) (*Engine, error) {
	invertedDB := storage.NewInvertedDB(
		termDB, forwardDB)
	return &Engine{
		invertedDB: invertedDB,
		bufSize:    1,
		N:          2,
	}, nil
}
