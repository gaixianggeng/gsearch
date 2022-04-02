package index

import "fmt"

// AddDocument 添加文档
func (e *Engine) AddDocument(title, body []byte) error {
	docID, err := e.db.AddDoc(title, body)
	if err != nil {
		return fmt.Errorf("AddDocument err: %v", err)
	}
	fmt.Println(docID)
	e.text2PostingsLists(docID, title)

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
