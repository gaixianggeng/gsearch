package query

import (
	"errors"
	"gsearch/pkg/utils/log"
	"strings"
)

// NGram 分词
func NGram(content string, n int32) ([]Tokenization, error) {
	if n < 1 {
		return nil, errors.New("ngram n must >= 1")
	}
	content = ignoredChar(content)
	var token []Tokenization
	if n >= int32(len([]rune(content))) {
		return []Tokenization{{content, 0}}, nil
	}

	i := int32(0)
	num := len([]rune(content))
	for i = 0; i < int32(num); i++ {
		if i+n > int32(num) {
			break
		}
		t := []rune(content)[i : i+n]
		token = append(token, Tokenization{
			Token:    string(t),
			Position: uint64(i),
		})
	}
	return token, nil
}

// ignoredChar 去除无用字符
func ignoredChar(str string) string {
	for _, c := range str {
		switch c {
		case ' ', '\f', '\n', '\r', '\t', '\v', '!', '"', '#', '$', '%', '&',
			'\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>',
			'?', '@', '[', '\\', ']', '^', '_', '`', '{', '|', '}', '~',
			0x3000, 0x3001, 0x3002, 0xFF08, 0xFF09, 0xFF01, 0xFF0C, 0xFF1A, 0xFF1B, 0xFF1F:
			str = strings.ReplaceAll(str, string(c), "")

		}
	}
	log.Debug(str)
	return str
}
