package recall

import (
	"brain/internal/engine"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

// Recall 查询召回
type Recall struct {
	*engine.Engine
	docCount     uint64 // 文档总数，用于计算相关性
	enablePhrase bool
	queryToken   []*queryTokenHash
}

// 用于实现排序map
type queryTokenHash struct {
	token         string
	invertedIndex *engine.InvertedIndexValue
}

// SearchItem 查询结果
type SearchItem struct {
	DocID uint64
	Score float64
}

// Recalls 召回结果
type Recalls []*SearchItem

// token游标 标识当前位置
type searchCursor struct {
	doc     *engine.PostingsList // 文档编号的序列
	current *engine.PostingsList // 当前的文档编号
}

// 短语游标
type phraseCursor struct {
	positions []uint64 // 位置信息
	base      uint64   // 词元在查询中的位置
	current   *uint64  // 当前的位置信息
	index     uint     // 当前的位置index
}

// Search 入口
func (r *Recall) Search(query string) (Recalls, error) {
	err := r.splitQuery2Tokens(query)
	if err != nil {
		log.Errorf("splitQuery2Tokens err: %v", err)
		return nil, fmt.Errorf("splitQuery2Tokens err: %v", err)
	}
	r.sortToken(r.Engine.PostingsHashBuf)
	if len(r.queryToken) == 0 {
		return nil, fmt.Errorf("queryTokenHash is nil")
	}
	return r.searchDoc()
}

func (r *Recall) splitQuery2Tokens(query string) error {
	err := r.Text2PostingsLists(query, 0)
	if err != nil {
		return fmt.Errorf("text2postingslists err: %v", err)
	}
	return nil
}

func (r *Recall) searchDoc() (Recalls, error) {

	recalls := make(Recalls, 0)

	tokenCount := len(r.queryToken)
	cursors := make([]searchCursor, tokenCount)

	// 为每个token初始化游标
	for i, t := range r.queryToken {
		// 正常不会出现，以防未知bug，所以设置fatal
		if t.token == "" {
			return nil, fmt.Errorf("token is nil")
		}
		postings, _, err := r.Engine.FetchPostings(t.token)
		if err != nil {
			return nil, fmt.Errorf("fetchPostings err: %v", err)
		}
		if postings == nil {
			return nil, fmt.Errorf("postings is nil")
		}

		log.Debugf("token:%v,invertedIndex:%v", t.token, postings.DocID)
		cursors[i].doc = postings
		cursors[i].current = postings
	}

	// 整个遍历token来匹配doc
	for cursors[0].current != nil {
		var docID, nextDocID uint64
		// 拥有文档最少的token作为标尺
		docID = cursors[0].current.DocID

		// 匹配其他token的doc
		for i := 1; i < tokenCount; i++ {
			cur := cursors[i]
			for cur.current != nil && cur.current.DocID < docID {
				cur.current = cur.current.Next
			}
			// 存在token关联的docid都小于cursors[0]的docID，则跳出
			if cur.current == nil {
				log.Infof("cur.current is nil\n")
				break
			}
			// 对于除词元A以外的词元，如果其document_id不等于词元A的document_id
			// 那么就将这个document_id设定为next_doc_id
			if cur.current.DocID != docID {
				nextDocID = cur.current.DocID
				break
			}
		}

		log.Debugf("当前docID:%v,nextDocID:%v", docID, nextDocID)
		if nextDocID > 0 {
			// 不断获取A的下一个document_id，直到其当前的document_id不小于next_doc_id为止
			for cursors[0].current != nil && cursors[0].current.DocID < nextDocID {
				cursors[0].current = cursors[0].current.Next
			}
		} else {
			// 有匹配上的docid
			phraseCount := int64(-1)
			if r.enablePhrase {
				phraseCount = r.searchPhrase(r.queryToken, cursors)
			}
			score := 0.0
			if phraseCount > 0 {
				// TODO: 计算相关性
				r.calculateScore(cursors, uint64(tokenCount))
			}
			cursors[0].current = cursors[0].current.Next
			log.Debugf("匹配召回docID:%v,nextDocID:%v,phrase:%d", docID, nextDocID, phraseCount)
			recalls = append(recalls, &SearchItem{DocID: docID, Score: score})
		}
	}
	log.Debugf("recalls size:%v", len(recalls))
	return recalls, nil
}

// 计算相关性
func (r *Recall) calculateScore(cursor []searchCursor, tokenCount uint64) float64 {
	return 0.0
}

// queryToken 查询query的倒排索引 tokenCursors是fetched文档的倒排索引
// 返回检索出的短语数
func (r *Recall) searchPhrase(queryToken []*queryTokenHash, tokenCursors []searchCursor) int64 {

	// 获取遍历查询query分词之后的词元总数 也就是被分成了多少个term
	positionsSum := uint64(0)
	for _, t := range queryToken {
		positionsSum += t.invertedIndex.PositionCount
	}

	cursors := make([]phraseCursor, positionsSum)
	phraseCount := int64(0)
	// 初始化游标 获取token关联的第一篇doc的pos相关数据
	n := 0
	for i, t := range queryToken {
		for _, pos := range t.invertedIndex.PostingsList.Positions {
			cursors[n].base = pos                                    // 记录查询中出现的位置
			cursors[n].positions = tokenCursors[i].current.Positions // 获取token的positions
			cursors[n].current = &cursors[i].positions[0]            // 获取文档中出现的位置
			cursors[n].index = 0                                     // 获取文档中出现的位置
			n++
		}
	}

	for cursors[0].current != nil {
		var relPos, nextRelPos uint64
		relPos = *cursors[0].current - cursors[0].base
		nextRelPos = relPos
		/* 对于除词元A以外的词元，不断地向后读取其出现位置，直到其偏移量不小于词元A的偏移量为止 */
		for i := 1; i < len(cursors); i++ {
			cur := &cursors[i]
			for cur.current != nil && *cur.current-cur.base < relPos {
				cur.index++
				if int(cur.index) >= len(cur.positions) {
					log.Warnf("cur.index >= len(cur.positions)\n")
					cur.current = nil
					break
				}
				cur.current = &cur.positions[cur.index]
			}
			if cur.current == nil {
				break
			}
			if *cur.current-cur.base != relPos {
				nextRelPos = *cur.current - cur.base
				break
			}
		}
		if nextRelPos > relPos {
			/* 不断向后读取，直到词元A的偏移量不小于next_rel_position为止 */
			for cursors[0].current != nil && *cursors[0].current-cursors[0].base < nextRelPos {
				cursors[0].index++
				if int(cursors[0].index) >= len(cursors[0].positions) {
					log.Warnf("cursors[0].index >= len(cursors[0].positions)\n")
					cursors[0].current = nil
					break
				}
				cursors[0].current = &cursors[0].positions[cursors[0].index]
			}
		} else {
			// 找到短语
			phraseCount++
			cursors[0].index++
			if int(cursors[0].index) >= len(cursors[0].positions) {
				log.Warnf("cursors[0].index >= len(cursors[0].positions)\n")
				cursors[0].current = nil
			} else {
				cursors[0].current = &cursors[0].positions[cursors[0].index]
			}
		}
	}

	return phraseCount
}

// token 根据doc count升序排序
func (r *Recall) sortToken(postHash engine.InvertedIndexHash) {
	tokenHash := make([]*queryTokenHash, 0)
	for token, invertedIndex := range postHash {
		q := new(queryTokenHash)
		q.token = token
		q.invertedIndex = invertedIndex
		tokenHash = append(tokenHash, q)
	}
	sort.Sort(docCountSort(tokenHash))
	r.queryToken = tokenHash
	for _, t := range r.queryToken {
		log.Debugf("token:%v,docCount:%v", t.token, t.invertedIndex.DocCount)
	}
	return
}

// NewRecall new
func NewRecall(e *engine.Engine) *Recall {
	docCount, err := e.ForwardDB.Count()
	if err != nil {
		log.Fatalf("e.ForwardDB.Count() error:%v\n", err)
		return nil
	}
	return &Recall{e, docCount, false, nil}
}
