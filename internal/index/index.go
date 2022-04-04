package index

import "fmt"

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
		e.indexCount++
	}

	// 落盘操作

	return nil

}

//
func createNewPostingList() *PostingsList {
	p := new(PostingsList)
	return p
}

// NewIndexEngine init
func NewIndexEngine() *Engine {
	return &Engine{}
}
