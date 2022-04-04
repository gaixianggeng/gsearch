package query

// Tokenization  分词返回结构
type Tokenization struct {
	Token    []byte
	Position int64
}

// Ngram 分词
func Ngram(content string, n int32) ([]Tokenization, error) {

	return nil, nil
}
