package engine

import (
	utils "brain/utils"
	"bytes"

	log "github.com/sirupsen/logrus"
)

// MergePostings merge two postings list
// https://leetcode-cn.com/problems/he-bing-liang-ge-pai-xu-de-lian-biao-lcof/
func MergePostings(pa, pb *PostingsList) *PostingsList {
	ret := new(PostingsList)
	p := new(PostingsList)
	p = nil
	for pa != nil || pb != nil {

		temp := new(PostingsList)
		if pb == nil || (pa != nil && pa.DocID <= pb.DocID) {
			temp = pa
			pa = pa.Next
		} else if pa == nil || (pb != nil && pa.DocID > pb.DocID) {
			temp = pb
			pb = pb.Next
		} else {
			break
		}
		temp.Next = nil

		if p == nil {
			ret.Next = temp
		} else {
			p.Next = temp
		}

		p = temp
	}

	return ret.Next
}

// MergeInvertedIndex 合并两个倒排索引
func MergeInvertedIndex(base, toBeAdded InvertedIndexHash) {
	for token, index := range base {
		if toBeAddedIndex, ok := (toBeAdded)[token]; ok {
			log.Debug("mergeInvertedIndex tokenID:", token)
			// 不需要+=positionCount 查询时候用到的字段，不需要写入到倒排表中
			index.PostingsList = MergePostings(index.PostingsList, toBeAddedIndex.PostingsList)
			index.DocsCount += toBeAddedIndex.DocsCount
			delete(toBeAdded, token)
		}
	}
	for tokenID, index := range toBeAdded {
		(base)[tokenID] = index
	}

}

// decodePostings 解码 return *PostingsList postingslen err
func decodePostings(buf *bytes.Buffer) (*PostingsList, uint64, error) {
	if buf == nil || buf.Len() == 0 {
		return nil, 0, nil
	}
	// p := new(PostingsList)
	// buf := bytes.NewBuffer([]byte{})
	// utils.BinaryRead(buf)

	return nil, 0, nil
}

// EncodePostings 编码
// bytes.Buffer
// docCount暂时用不到
func EncodePostings(postings *PostingsList, postingsLen uint64) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := utils.BinaryWrite(buf, postingsLen)
	if err != nil {
		return nil, err
	}

	log.Debug(len(buf.Bytes()))
	for postings != nil {
		log.Debugf("docid:%d,count:%d,positions:%v", postings.DocID, postings.PositionCount, postings.Positions)
		err := utils.BinaryWrite(buf, postings.DocID)
		if err != nil {
			return nil, err
		}
		err = utils.BinaryWrite(buf, postings.PositionCount)
		if err != nil {
			return nil, err
		}
		err = utils.BinaryWrite(buf, postings.Positions)
		if err != nil {
			return nil, err
		}
		postings = postings.Next
	}
	log.Debug(len(buf.Bytes()))
	return buf, nil
}
