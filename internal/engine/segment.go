package engine

import (
	"bytes"
	"doraemon/internal/storage"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// SegInfo 段信息
type SegInfo struct {
	SegID            uint64 `json:"seg_name"`           // 段前缀名
	SegSize          uint64 `json:"seg_size"`           // 写入doc数量
	InvertedFileSize uint64 `json:"inverted_file_size"` // 写入inverted文件大小
	ForwardFileSize  uint64 `json:"forward_file_size"`  // 写入forward文件大小
	DelSize          uint64 `json:"del_size"`           // 删除文档数量
	DelFileSize      uint64 `json:"del_file_size"`      // 删除文档文件大小
	TermSize         uint64 `json:"term_size"`          // term文档文件大小
	TermFileSize     uint64 `json:"term_file_size"`     // term文件大小
	ReferenceCount   uint64 `json:"reference_count"`    // 引用计数
	IsReading        bool   `json:"is_reading"`         // 是否正在被读取
	IsMerging        bool   `json:"is_merging"`         // 是否正在参与合并
}

// https://www.cnblogs.com/qianye/archive/2012/11/25/2787923.html

// LoserTree --
type LoserTree struct {
	tree   []int // 索引表示顺序，0表示最小值，value表示对应的leaves的index
	leaves []*TermNode
}

// TermNode --
type TermNode struct {
	Info  chan storage.TermInfo
	Key   []byte
	Value []byte
	DB    *storage.InvertedDB
}

// NewSegLoserTree 败者树
func NewSegLoserTree(leaves []*TermNode) *LoserTree {
	k := len(leaves)

	lt := &LoserTree{
		tree:   make([]int, k),
		leaves: leaves,
	}
	if k > 0 {
		lt.initWinner(0)
	}
	return lt
}

// 整体逻辑 输的留下来，赢的向上比
// 获取最小值的索引
func (lt *LoserTree) initWinner(idx int) int {
	log.Debugf("idx:%d", idx)
	// 根节点有一个父节点，存储最小值。
	if idx == 0 {
		lt.tree[0] = lt.initWinner(1)
		return lt.tree[0]
	}
	if idx >= len(lt.tree) {
		return idx - len(lt.tree)
	}

	left := lt.initWinner(idx * 2)
	right := lt.initWinner(idx*2 + 1)
	log.Debugf("left:%d, right:%d", left, right)

	// 左不为空，右为空，则记录右边
	if lt.leaves[left] != nil && lt.leaves[right] == nil {
		left, right = right, left

	}
	if lt.leaves[left] != nil && lt.leaves[right] != nil {
		leftCh := <-lt.leaves[left].Info
		rightCh := <-lt.leaves[right].Info
		lt.leaves[left].Key = leftCh.Key
		lt.leaves[left].Value = leftCh.Value
		lt.leaves[right].Key = rightCh.Key
		lt.leaves[right].Value = rightCh.Value

		log.Debugf("leftCh:%s, rightCh:%s", leftCh.Key, rightCh.Key)
		if string(leftCh.Key) < string(rightCh.Key) {
			left, right = right, left
		}

	}
	// 左边的节点比右边的节点大
	// 记录败者 即 记录较大的节点索引 较小的继续向上比较
	lt.tree[idx] = left
	return right
}

// Pop 弹出最小值
func (lt *LoserTree) Pop() *TermNode {
	if len(lt.tree) == 0 {
		return nil
	}

	// 取出最小的索引
	leafWinIdx := lt.tree[0]
	// 找到对应叶子节点
	winner := lt.leaves[leafWinIdx]

	target := new(TermNode)

	// 更新对应index里节点的值
	// 如果是最后一个节点，标识nil
	if winner == nil {
		log.Debugf("数据已读取完毕 winner.Key == nil")
		lt.leaves[leafWinIdx] = nil
		target = nil
	} else {
		log.Debugf("winner:%s", winner.Key)
		// 赋值
		target.Key = winner.Key
		target.Value = winner.Value
		target.DB = winner.DB

		// 获取下一轮的key和value
		termCh, isOpen := <-winner.Info
		// channel已关闭
		if !isOpen {
			log.Debugf("channel已关闭")
			lt.leaves[leafWinIdx] = nil
		} else {
			// 重新赋值
			lt.leaves[leafWinIdx].Key = termCh.Key
			lt.leaves[leafWinIdx].Value = termCh.Value
		}

	}

	// 获取父节点
	treeIdx := (leafWinIdx + len(lt.tree)) / 2

	for treeIdx != 0 {
		loserLeafIdx := lt.tree[treeIdx]
		// 如果为nil，则将父节点的idx设置为该索引，不为空的继续向上比较
		if lt.leaves[loserLeafIdx] == nil {
			lt.tree[treeIdx] = loserLeafIdx
		} else {
			// 如果已经该叶子节点已经读取完毕，则将父节点的idx设置为该索引
			if lt.leaves[leafWinIdx] == nil {
				loserLeafIdx, leafWinIdx = leafWinIdx, loserLeafIdx
			} else if string(lt.leaves[loserLeafIdx].Key) <
				string(lt.leaves[leafWinIdx].Key) {
				loserLeafIdx, leafWinIdx = leafWinIdx, loserLeafIdx
			}
			// 更新
			lt.tree[treeIdx] = loserLeafIdx
		}
		treeIdx /= 2
	}
	lt.tree[0] = leafWinIdx

	time.Sleep(1e8)

	return target
}

// MergeKSegments 多路归并
func MergeKSegments(lists []*TermNode) InvertedIndexHash {
	// var dummy = &ListNode{}
	// pre := dummy
	log.Debugf("start merge k segemnts[lists:%v]", lists)
	lt := NewSegLoserTree(lists)
	log.Debugf("init:%s,%s", string(lt.leaves[0].Key), string(lt.leaves[1].Key))

	log.Debugf("init:%+v", lt)

	log.Debugf(strings.Repeat("-", 20))
	res := make(InvertedIndexHash)
	for {
		node := lt.Pop()
		if node == nil {
			break
		}
		val, err := node.DB.Bytes2TermVal(node.Value)
		if err != nil {
			log.Errorf("bytes2termval err:%s", err)
			continue
		}
		// 解码
		c, err := node.DB.GetDocInfo(val[1], val[2])
		if err != nil {
			log.Errorf("FetchPostings getDocInfo err: %v", err)
			continue
		}
		pos, count, err := decodePostings(bytes.NewBuffer(c))

		log.Debugf("pop node key:%+v,value:%v,count:%d", string(node.Key), val, count)
		log.Debugf(strings.Repeat("-", 20))

		if p, ok := res[string(node.Key)]; ok {
			p.DocCount += count
			p.PostingsList = MergePostings(p.PostingsList, pos)
			continue
		}
		res[string(node.Key)] = &InvertedIndexValue{
			Token:        string(node.Key),
			DocCount:     count,
			PostingsList: pos,
		}
		// pre.Next = node
		// pre = node
	}
	return res
}
