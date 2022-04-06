package index

import (
	"fmt"
)

// AddDocument 添加文档
func (e *Engine) AddDocument(title, body []byte) error {
	if len(title) > 0 || len(body) > 0 {
		docID, err := e.db.Add(title, body)
		if err != nil {
			return fmt.Errorf("AddDocument err: %v", err)
		}
		fmt.Println(docID)
		err = e.text2PostingsLists(docID, title)
		if err != nil {
			return fmt.Errorf("text2postingslists err: %v", err)
		}
		e.bufCount++

		e.indexCount++
	}

	// 落盘操作
	if len(e.postingsHashBuf) > 0 && (e.bufCount > e.bufSize || title == nil) {

		for tokenID, invertedIndex := range e.postingsHashBuf {

			fmt.Printf("tokenID:%d,invertedIndex:%v\n", tokenID, invertedIndex)
			updatePostings(invertedIndex)
		}

		// 重置
		e.postingsHashBuf = make(InvertedIndexHash)
		e.bufCount = 0
	}

	return nil

}

// 创建倒排列表
func createNewPostingList(docID int64) *PostingsList {
	p := new(PostingsList)
	p.DocID = docID
	p.positionCount = 1
	p.positions = make([]int64, 0)
	return p
}

// 创建倒排索引
func createNewInvertedIndex(tokenID, docCount int64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.docsCount = docCount
	p.TokenID = tokenID
	p.positionCount = 0
	p.postingList = new(PostingsList)
	return p
}

// NewIndexEngine init
func NewIndexEngine() *Engine {
	return &Engine{
		bufSize: 30,
	}
}
