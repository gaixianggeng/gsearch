package query

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Tokenization  分词返回结构
type Tokenization struct {
	Token    string
	Position uint64
}

// Ngram 分词
func Ngram(content string, n int32) ([]Tokenization, error) {
	if n < 1 {
		return nil, fmt.Errorf("Ngram n must >= 1")
	}
	content = ignoredChar(content)
	var token []Tokenization
	if n >= int32(len([]rune(content))) {
		token = append(token, Tokenization{content, 0})
		return token, nil
	}

	i := int32(0)
	num := len([]rune(content))
	for i = 0; i < int32(num); i++ {
		t := []rune{}
		if i+n > int32(num) {
			// t = []rune(content)[i:]
			break
		} else {
			t = []rune(content)[i : i+n]
		}
		log.Infof("tokens: %v", string(t))
		token = append(token, Tokenization{
			Token:    string(t),
			Position: uint64(i),
		})
	}
	return token, nil
}

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
