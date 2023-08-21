package tests

import (
	"gsearch/pkg/utils/log"
	"testing"
)

// 作者：ppd-2
// 链接：https://leetcode-cn.com/problems/merge-k-sorted-lists/solution/xcdytv-by-ppd-2-sspr/
// 来源：力扣（LeetCode）
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

// https://www.cnblogs.com/qianye/archive/2012/11/25/2787923.html
type ListNode struct {
	Val  int
	Next *ListNode
}

type LoserTree struct {
	tree   []int // 索引表示顺序，0表示最小值，value表示对应的leaves的index
	leaves []*ListNode
}

// var dummyVal = 100000
// var dummyListNode = ListNode{Val: dummyVal}

func NewLoserTreee(leaves []*ListNode) *LoserTree {
	k := len(leaves)
	// // 奇数+1
	// if k&1 == 1 {
	// 	leaves = append(leaves, &dummyListNode)
	// 	k++
	// }
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

	// // 为空的添加一个最大值
	// if lt.leaves[left] == nil {
	// 	lt.leaves[left] = dummyListNode
	// }
	// if lt.leaves[right] == nil {
	// 	lt.leaves[right] = &dummyListNode
	// }

	if lt.leaves[left].Val < lt.leaves[right].Val {
		left, right = right, left
	}
	// 左边的节点比右边的节点大
	// 记录败者 即 记录较大的节点索引 较小的继续向上比较
	lt.tree[idx] = left
	return right
}

func (lt *LoserTree) Pop() *ListNode {
	if len(lt.tree) == 0 {
		// return &dummyListNode
		return nil
	}

	// 取出最小的索引
	leafWinIdx := lt.tree[0]
	// 找到对应叶子节点
	winner := lt.leaves[leafWinIdx]

	// 更新对应index里节点的值
	if winner == nil {
		lt.leaves[leafWinIdx] = nil
	} else {
		lt.leaves[leafWinIdx] = winner.Next
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
			} else if lt.leaves[loserLeafIdx].Val <
				lt.leaves[leafWinIdx].Val {
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

func mergeKLists(lists []*ListNode) *ListNode {
	var dummy = &ListNode{}
	pre := dummy
	lt := NewLoserTreee(lists)

	log.Debugf("init:%+v", lt)
	for {
		node := lt.Pop()
		// if node == &dummyListNode {
		// 	break
		// }
		if node == nil {
			break
		}
		pre.Next = node
		pre = node
	}
	return dummy
}

func TestTree(t *testing.T) {
	// [[1,4,5],[3,3,4],[2,6]]
	listNode1 := &ListNode{Val: 1}
	listNode1.Next = &ListNode{Val: 4}
	listNode1.Next.Next = &ListNode{Val: 5}

	listNode2 := &ListNode{Val: 3}
	listNode2.Next = &ListNode{Val: 3}
	listNode2.Next.Next = &ListNode{Val: 4}

	listNode3 := &ListNode{Val: 2}
	listNode3.Next = &ListNode{Val: 6}

	list := mergeKLists([]*ListNode{listNode1, listNode2, listNode3})
	crossListNode(list)
}

func crossListNode(list *ListNode) {
	for list != nil {
		log.Debugf("%d", list.Val)
		list = list.Next
	}

}
