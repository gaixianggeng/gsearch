package analyzer

// StopWords is a set of common words to be ignored in text analysis
var StopWords = map[string]bool{
	// TODO: Add common stop words
}

// IsStopWord checks if a word is a stop word
func IsStopWord(word string) bool {
	_, exists := StopWords[word]
	return exists
}
