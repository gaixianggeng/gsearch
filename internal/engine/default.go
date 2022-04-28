package engine

var (
	// TermDBSuffix term db suffix
	TermDBSuffix = ".term"
	// InvertedDBSuffix inverted db suffix
	InvertedDBSuffix = ".inverted"
	// ForwardDBSuffix forward db suffix
	ForwardDBSuffix = ".forward"
)
var (
	termName     = ""
	invertedName = ""
	forwardName  = ""
)

// Mode 查询or索引模式
type Mode int32

const (
	// SearchMode 查询模式
	SearchMode Mode = 1
	// IndexMode 索引模式
	IndexMode Mode = 2
	// MergeMode seg merge模式
	MergeMode Mode = 3
)
