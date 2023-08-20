package storage

// Document 文档格式
type Document struct {
	DocID uint64
	Title string
	Body  string
}

// KvInfo term信息
type KvInfo struct {
	Key   []byte
	Value []byte
}

const forwardBucket = "forward"

const ForwardCountKey = "forwardCount"
