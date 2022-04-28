package engine

import (
	"doraemon/internal/storage"

	"github.com/boltdb/bolt"
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
	*bolt.Cursor
	*storage.TermInfo
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

	if string(lt.leaves[left].Key) < string(lt.leaves[right].Key) {
		left, right = right, left
	}
	// 左边的节点比右边的节点大
	// 记录败者 即 记录较大的节点索引 较小的继续向上比较
	lt.tree[idx] = left
	return right
}

// Pop 弹出最小值
func (lt *LoserTree) Pop() *TermNode {
	if len(lt.tree) == 0 {
		// return &dummyListNode
		return nil
	}

	// 取出最小的索引
	leafWinIdx := lt.tree[0]
	// 找到对应叶子节点
	winner := lt.leaves[leafWinIdx]

	// 更新对应index里节点的值
	// 如果是最后一个节点，标识nil
	if winner == nil {
		lt.leaves[leafWinIdx] = nil
	} else {
		key, value := winner.Next()
		lt.leaves[leafWinIdx] = &TermNode{
			Cursor: winner.Cursor,
			TermInfo: &storage.TermInfo{
				Key:   key,
				Value: value,
			},
		}
	}

	// if lt.leaves[leafWinIdx] == nil {
	// 	lt.leaves[leafWinIdx] = &dummyListNode
	// }

	// 获取父节点
	treeIdx := (leafWinIdx + len(lt.tree)) / 2

	log.Debugf("treeIdx:%d, leafWinIdx:%d", treeIdx, leafWinIdx)

	for treeIdx != 0 {
		// 如果第二小的节点比新取出的叶子节点的值小，则互换位置
		loserLeafIdx := lt.tree[treeIdx]
		if lt.leaves[loserLeafIdx] == nil {
			lt.tree[treeIdx] = loserLeafIdx
		} else {
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
	return winner
}

// MergeKSegments 多路归并
func MergeKSegments(lists []*TermNode) {
	// var dummy = &ListNode{}
	// pre := dummy
	log.Debugf("start merge k segemnts[lists:%v]", lists)
	lt := NewSegLoserTree(lists)

	log.Debugf("init:%+v", lt)
	for {
		node := lt.Pop()
		if node == nil {
			break
		}
		log.Debugf("node:%+v", string(node.Key))
		// pre.Next = node
		// pre = node
	}
	return
}
