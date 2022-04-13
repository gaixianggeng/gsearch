package index

import (
	"brain/internal/query"
	"bytes"
	"encoding/binary"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// mergePostings merge two postings list
// https://leetcode-cn.com/problems/he-bing-liang-ge-pai-xu-de-lian-biao-lcof/
func mergePostings(pa, pb *PostingsList) *PostingsList {
	ret := new(PostingsList)
	p := new(PostingsList)
	p = nil
	for pa != nil || pb != nil {

		temp := new(PostingsList)
		if pb == nil || (pa != nil && pa.DocID <= pb.DocID) {
			temp = pa
			pa = pa.next
		} else if pa == nil || (pb != nil && pa.DocID > pb.DocID) {
			temp = pb
			pb = pb.next
		} else {
			break
		}
		temp.next = nil

		if p == nil {
			ret.next = temp
		} else {
			p.next = temp
		}

		p = temp
	}

	return ret.next
}

// mergeInvertedIndex 合并两个倒排索引
func mergeInvertedIndex(base, toBeAdded InvertedIndexHash) {
	for tokenID, index := range base {
		if toBeAddedIndex, ok := (toBeAdded)[tokenID]; ok {
			index.postingsList = mergePostings(index.postingsList, toBeAddedIndex.postingsList)
			index.docsCount += toBeAddedIndex.docsCount
			delete(toBeAdded, tokenID)
		}
	}
	for tokenID, index := range toBeAdded {
		(base)[tokenID] = index
	}

}

// 解码
func decodePostings() {

}

// 编码
// bytes.Buffer
// docCount暂时用不到
func encodePostings(postings *PostingsList, postingsLen uint64) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := binaryWrite(buf, postingsLen)
	if err != nil {
		return nil, err
	}
	for postings != nil {
		err := binaryWrite(buf, postings.DocID)
		if err != nil {
			return nil, err
		}
		err = binaryWrite(buf, postings.positionCount)
		if err != nil {
			return nil, err
		}
		err = binaryWrite(buf, postings.positions)
		if err != nil {
			return nil, err
		}
		postings = postings.next
	}
	// binary.Write(buf, binary.BigEndian, postingsLen)
	return buf, nil
}

func fetchPostings(token string) (*PostingsList, uint64, error) {

	return nil, 0, nil
}

func (e *Engine) updatePostings(p *InvertedIndexValue) error {
	if p == nil {
		return fmt.Errorf("updatePostings p is nil")
	}
	// 拉取数据库数据
	oldPostings, size, err := fetchPostings(p.Token)
	if err != nil {
		return fmt.Errorf("updatePostings fetchPostings err: %v", err)
	}
	// merge
	if size > 0 {
		p.postingsList = mergePostings(oldPostings, p.postingsList)
		p.docsCount += size
	}
	// 开始写入数据库
	buf, err := encodePostings(p.postingsList, p.docsCount)
	if err != nil {
		return fmt.Errorf("updatePostings encodePostings err: %v", err)
	}
	return e.invertedDB.DBUpdatePostings(p.Token, buf.Bytes())
}

// text2PostingsLists --
func (e *Engine) text2PostingsLists(docID uint64, text []byte) error {
	tokens, err := query.Ngram(string(text), e.N)
	if err != nil {
		return fmt.Errorf("text2PostingsLists Ngram err: %v", err)
	}
	bufInvertedHash := make(InvertedIndexHash)

	for _, token := range tokens {
		err := e.token2PostingsLists(bufInvertedHash, token.Token, token.Position, docID)
		if err != nil {
			return fmt.Errorf("text2PostingsLists token2PostingsLists err: %v", err)
		}
	}

	if e.postingsHashBuf != nil && len(e.postingsHashBuf) > 0 {
		mergeInvertedIndex(e.postingsHashBuf, bufInvertedHash)
	} else {
		e.postingsHashBuf = make(InvertedIndexHash)
		e.postingsHashBuf = bufInvertedHash
	}
	return nil

}

func (e *Engine) token2PostingsLists(bufInvertHash InvertedIndexHash, token string,
	position uint64, docID uint64) error {

	// init
	bufInvert := new(InvertedIndexValue)

	// doc_id用来标识写入数据还是查询数据
	docCount, err := e.tokenDB.GetToken(token, docID)
	if err != nil {
		return fmt.Errorf("token2PostingsLists GetTokenID err: %v", err)
	}

	if len(bufInvertHash) > 0 {
		if item, ok := bufInvertHash[token]; ok {
			bufInvert = item
		}
	}

	pl := new(PostingsList)
	if bufInvert != nil && bufInvert.postingsList != nil {
		log.Debug("token2PostingsLists bufInvert.postingsList is not nil")
		pl = bufInvert.postingsList
		// 这里的positioinCount和下面bufInvert的positionCount是不一样的
		// 这里统计的是同一个docid的position的个数
		pl.positionCount++
	} else {

		log.Debug("token2PostingsLists bufInvert.postingsList is nil")
		if docID != 0 {
			docCount = 1
		}
		bufInvert = createNewInvertedIndex(token, docCount)
		bufInvertHash[token] = bufInvert
		pl = createNewPostingsList(docID)
		bufInvert.postingsList = pl
	}
	// 存储位置信息
	pl.positions = append(pl.positions, position)
	// 统计该token关联的所有doc的position的个数
	bufInvert.positionCount++

	return nil
}

func binaryWrite(buf *bytes.Buffer, v any) error {
	size := binary.Size(v)
	log.Debug("docid size:", size)
	if size <= 0 {
		return fmt.Errorf("encodePostings binary.Size err,size: %v", size)
	}
	return binary.Write(buf, binary.LittleEndian, v)
}
