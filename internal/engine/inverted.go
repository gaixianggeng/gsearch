package engine

//InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string
	PostingsList  *PostingsList
	DocCount      uint64
	PositionCount uint64 // 查询使用，写入的时候暂时用不到
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*InvertedIndexValue

// CreateNewInvertedIndex 创建倒排索引
func CreateNewInvertedIndex(token string, docCount uint64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.DocCount = docCount
	p.Token = token
	p.PositionCount = 0
	p.PostingsList = new(PostingsList)
	return p
}
