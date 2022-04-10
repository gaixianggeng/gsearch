package bptree

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"syscall"
)

var (
	err   error
	order = 4
)

const (
	// InvalidOffset --
	InvalidOffset = 0xdeadbeef
	// MaxFreeBlocks --
	MaxFreeBlocks = 100
)

// ErrorHasExistedKey --
var ErrorHasExistedKey = errors.New("hasExistedKey")

// ErrorNotFoundKey --
var ErrorNotFoundKey = errors.New("notFoundKey")

// ErrorInvalidDBFormat --
var ErrorInvalidDBFormat = errors.New("invalid db format")

// OFFTYPE --
type OFFTYPE uint64

// Tree b+树自身信息
type Tree struct {
	rootOff    OFFTYPE
	nodePool   *sync.Pool
	freeBlocks []OFFTYPE
	file       *os.File
	blockSize  uint64 // block 大小
	fileSize   uint64 // b+树文件大小
}

// Node 节点信息
// key和childern是等价的 key存储的是值，childern存储的是子节点的offset
type Node struct {
	IsActive bool // 节点所在的磁盘空间是否在当前b+树内
	Children []OFFTYPE
	Self     OFFTYPE
	Next     OFFTYPE
	Prev     OFFTYPE
	Parent   OFFTYPE
	Keys     []uint64
	Values   [][]byte
	IsLeaf   bool // 是否为叶子节点
}

// NewTree --
func NewTree(filename string) (*Tree, error) {
	var (
		stat  syscall.Statfs_t
		fstat os.FileInfo
		err   error
	)

	t := &Tree{}

	t.rootOff = InvalidOffset
	t.nodePool = &sync.Pool{
		New: func() interface{} {
			return &Node{}
		},
	}
	t.freeBlocks = make([]OFFTYPE, 0, MaxFreeBlocks)
	if t.file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return nil, err
	}

	if err = syscall.Statfs(filename, &stat); err != nil {
		return nil, err
	}
	t.blockSize = uint64(stat.Bsize)
	if t.blockSize == 0 {
		return nil, errors.New("blockSize should be zero")
	}
	if fstat, err = t.file.Stat(); err != nil {
		return nil, err
	}

	t.fileSize = uint64(fstat.Size())
	if t.fileSize != 0 {
		if err = t.restructRootNode(); err != nil {
			return nil, err
		}
		if err = t.checkDiskBlockForFreeNodeList(); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// Close --
func (t *Tree) Close() error {
	if t.file != nil {
		t.file.Sync()
		return t.file.Close()
	}

	return nil
}

// 找根节点的offset
func (t *Tree) restructRootNode() error {
	var (
		err error
	)
	node := &Node{}

	//
	for off := uint64(0); off < t.fileSize; off += t.blockSize {
		if err = t.seekNode(node, OFFTYPE(off)); err != nil {
			return err
		}
		if node.IsActive {
			break
		}
	}
	if !node.IsActive {
		return ErrorInvalidDBFormat
	}
	for node.Parent != InvalidOffset {
		if err = t.seekNode(node, node.Parent); err != nil {
			return err
		}
	}

	t.rootOff = node.Self

	return nil
}

func (t *Tree) checkDiskBlockForFreeNodeList() error {
	var (
		err error
	)
	node := &Node{}
	bs := t.blockSize
	// 统计free blocks
	for off := uint64(0); off < t.fileSize && len(t.freeBlocks) < MaxFreeBlocks; off += bs {
		if off+bs > t.fileSize {
			break
		}
		if err = t.seekNode(node, OFFTYPE(off)); err != nil {
			return err
		}
		if !node.IsActive {
			t.freeBlocks = append(t.freeBlocks, OFFTYPE(off))
		}
	}

	// 分配满blocks
	nextFile := ((t.fileSize + 4095) / 4096) * 4096
	for len(t.freeBlocks) < MaxFreeBlocks {
		t.freeBlocks = append(t.freeBlocks, OFFTYPE(nextFile))
		nextFile += bs
	}
	t.fileSize = nextFile
	return nil
}

func (t *Tree) initNodeForUsage(node *Node) {
	node.IsActive = true
	node.Children = nil
	node.Self = InvalidOffset
	node.Next = InvalidOffset
	node.Prev = InvalidOffset
	node.Parent = InvalidOffset
	node.Keys = nil
	node.Values = nil
	node.IsLeaf = false
}

func (t *Tree) clearNodeForUsage(node *Node) {
	node.IsActive = false
	node.Children = nil
	node.Self = InvalidOffset
	node.Next = InvalidOffset
	node.Prev = InvalidOffset
	node.Parent = InvalidOffset
	node.Keys = nil
	node.Values = nil
	node.IsLeaf = false
}

// 根据off 读取node的相关信息
func (t *Tree) seekNode(node *Node, off OFFTYPE) error {
	if node == nil {
		return fmt.Errorf("cant use nil for seekNode")
	}
	t.clearNodeForUsage(node)

	var err error
	// 1. 从off读取8字节数据
	buf := make([]byte, 8)
	if n, err := t.file.ReadAt(buf, int64(off)); err != nil {
		return err
	} else if uint64(n) != 8 {
		return fmt.Errorf("readat %d from %s, expected len = %d but get %d", off, t.file.Name(), 4, n)
	}
	bs := bytes.NewBuffer(buf)

	// 反序列化 转10进制
	dataLen := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &dataLen); err != nil {
		return err
	}
	if uint64(dataLen)+8 > t.blockSize {
		return fmt.Errorf("flushNode len(node) = %d exceed t.blockSize %d", uint64(dataLen)+4, t.blockSize)
	}

	// 2. 继续读取其他字段相关信息
	buf = make([]byte, dataLen)
	if n, err := t.file.ReadAt(buf, int64(off)+8); err != nil {
		return err
	} else if uint64(n) != uint64(dataLen) {
		return fmt.Errorf("readat %d from %s, expected len = %d but get %d", uint64(off)+4, t.file.Name(), dataLen, n)
	}

	bs = bytes.NewBuffer(buf)
	// IsActive
	if err = binary.Read(bs, binary.LittleEndian, &node.IsActive); err != nil {
		return err
	}
	// Children
	childCount := uint8(0)
	if err = binary.Read(bs, binary.LittleEndian, &childCount); err != nil {
		return err
	}

	node.Children = make([]OFFTYPE, childCount)
	for i := uint8(0); i < childCount; i++ {
		child := uint64(0)
		if err = binary.Read(bs, binary.LittleEndian, &child); err != nil {
			return err
		}
		node.Children[i] = OFFTYPE(child)
	}
	// Self
	self := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &self); err != nil {
		return err
	}
	node.Self = OFFTYPE(self)
	// Next
	next := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &next); err != nil {
		return err
	}
	node.Next = OFFTYPE(next)
	// Prev
	prev := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &prev); err != nil {
		return err
	}
	node.Prev = OFFTYPE(prev)
	// Parent
	parent := uint64(0)
	if err = binary.Read(bs, binary.LittleEndian, &parent); err != nil {
		return err
	}
	node.Parent = OFFTYPE(parent)
	// Keys
	keysCount := uint8(0)
	if err = binary.Read(bs, binary.LittleEndian, &keysCount); err != nil {
		return err
	}
	node.Keys = make([]uint64, keysCount)
	for i := uint8(0); i < keysCount; i++ {
		if err = binary.Read(bs, binary.LittleEndian, &node.Keys[i]); err != nil {
			return err
		}
	}
	// Records
	recordCount := uint8(0)
	if err = binary.Read(bs, binary.LittleEndian, &recordCount); err != nil {
		return err
	}
	node.Values = make([][]byte, recordCount)
	for i := uint8(0); i < recordCount; i++ {
		l := uint8(0)
		if err = binary.Read(bs, binary.LittleEndian, &l); err != nil {
			return err
		}
		v := make([]byte, l)
		if err = binary.Read(bs, binary.LittleEndian, &v); err != nil {
			return err
		}
		node.Values[i] = v
	}
	// IsLeaf
	if err = binary.Read(bs, binary.LittleEndian, &node.IsLeaf); err != nil {
		return err
	}

	return nil
}

func (t *Tree) flushNodesAndPutNodesPool(nodes ...*Node) error {
	for _, n := range nodes {
		if err := t.flushNodeAndPutNodePool(n); err != nil {
			return err
		}
	}
	return err
}

func (t *Tree) flushNodeAndPutNodePool(n *Node) error {
	if err := t.flushNode(n); err != nil {
		return err
	}
	t.putNodePool(n)
	return nil
}

// 归还至pool
func (t *Tree) putNodePool(n *Node) {
	t.nodePool.Put(n)
}

// 将node写入文件
func (t *Tree) flushNode(n *Node) error {
	if n == nil {
		return fmt.Errorf("flushNode == nil")
	}
	if t.file == nil {
		return fmt.Errorf("flush node into disk, but not open file")
	}

	var (
		length int
		err    error
	)

	bs := bytes.NewBuffer(make([]byte, 0))

	// IsActive
	if err = binary.Write(bs, binary.LittleEndian, n.IsActive); err != nil {
		return nil
	}

	// Children
	childCount := uint8(len(n.Children))
	if err = binary.Write(bs, binary.LittleEndian, childCount); err != nil {
		return err
	}
	for _, v := range n.Children {
		if err = binary.Write(bs, binary.LittleEndian, uint64(v)); err != nil {
			return err
		}
	}

	// Self
	if err = binary.Write(bs, binary.LittleEndian, uint64(n.Self)); err != nil {
		return err
	}

	// Next
	if err = binary.Write(bs, binary.LittleEndian, uint64(n.Next)); err != nil {
		return err
	}

	// Prev
	if err = binary.Write(bs, binary.LittleEndian, uint64(n.Prev)); err != nil {
		return err
	}

	// Parent
	if err = binary.Write(bs, binary.LittleEndian, uint64(n.Parent)); err != nil {
		return err
	}

	// Keys
	keysCount := uint8(len(n.Keys))
	if err = binary.Write(bs, binary.LittleEndian, keysCount); err != nil {
		return err
	}
	for _, v := range n.Keys {
		if err = binary.Write(bs, binary.LittleEndian, v); err != nil {
			return err
		}
	}

	// Record
	recordCount := uint8(len(n.Values))
	if err = binary.Write(bs, binary.LittleEndian, recordCount); err != nil {
		return err
	}
	for _, v := range n.Values {
		if err = binary.Write(bs, binary.LittleEndian, uint8(len([]byte(v)))); err != nil {
			return err
		}
		if err = binary.Write(bs, binary.LittleEndian, []byte(v)); err != nil {
			return err
		}
	}

	// IsLeaf
	if err = binary.Write(bs, binary.LittleEndian, n.IsLeaf); err != nil {
		return err
	}

	dataLen := len(bs.Bytes())
	if uint64(dataLen)+8 > t.blockSize {
		return fmt.Errorf("flushNode len(node) = %d exceed t.blockSize %d", uint64(dataLen)+4, t.blockSize)
	}
	tmpbs := bytes.NewBuffer(make([]byte, 0))
	if err = binary.Write(tmpbs, binary.LittleEndian, uint64(dataLen)); err != nil {
		return err
	}

	data := append(tmpbs.Bytes(), bs.Bytes()...)
	if length, err = t.file.WriteAt(data, int64(n.Self)); err != nil {
		return err
	} else if len(data) != length {
		return fmt.Errorf("writeat %d into %s, expected len = %d but get %d", uint64(n.Self), t.file.Name(), len(data), length)
	}
	return nil
}

// 根据offset组装node
func (t *Tree) newMappingNodeFromPool(off OFFTYPE) (*Node, error) {
	node := t.nodePool.Get().(*Node)
	t.initNodeForUsage(node)
	if off == InvalidOffset {
		return node, nil
	}
	// 如果offset存在的话，说明是从文件中读取的，读取数据赋值给node
	t.clearNodeForUsage(node)
	if err := t.seekNode(node, off); err != nil {
		return nil, err
	}
	return node, nil
}

// 为node赋值磁盘的offset
func (t *Tree) newNodeFromDisk() (*Node, error) {
	var (
		node *Node
		err  error
	)
	node = t.nodePool.Get().(*Node)
	if len(t.freeBlocks) > 0 {
		off := t.freeBlocks[0]
		t.freeBlocks = t.freeBlocks[1:len(t.freeBlocks)]
		t.initNodeForUsage(node)
		node.Self = off
		return node, nil
	}
	if err = t.checkDiskBlockForFreeNodeList(); err != nil {
		return nil, err
	}
	if len(t.freeBlocks) > 0 {
		off := t.freeBlocks[0]
		t.freeBlocks = t.freeBlocks[1:len(t.freeBlocks)]
		t.initNodeForUsage(node)
		node.Self = off
		return node, nil
	}
	return nil, fmt.Errorf("can't not alloc more node")
}

// 归还offset
func (t *Tree) putFreeBlocks(off OFFTYPE) {
	if len(t.freeBlocks) >= MaxFreeBlocks {
		return
	}
	t.freeBlocks = append(t.freeBlocks, off)
}

// Find 查找node
func (t *Tree) Find(key uint64) ([]byte, error) {
	var (
		node *Node
		err  error
	)

	if t.rootOff == InvalidOffset {
		return []byte{}, nil
	}

	// 新建节点
	if node, err = t.newMappingNodeFromPool(InvalidOffset); err != nil {
		return []byte{}, err
	}

	// 找到key所在的node
	if err = t.findLeaf(node, key); err != nil {
		return []byte{}, err
	}
	defer t.putNodePool(node)

	// 匹配
	// TODO: 感觉可以用二分查找
	for i, nkey := range node.Keys {
		if nkey == key {
			return node.Values[i], nil
		}
	}
	return []byte{}, ErrorNotFoundKey
}

func (t *Tree) findLeaf(node *Node, key uint64) error {
	var (
		err  error
		root *Node
	)

	c := t.rootOff
	if c == InvalidOffset {
		return nil
	}

	// 查找根节点
	if root, err = t.newMappingNodeFromPool(c); err != nil {
		return err
	}
	defer t.putNodePool(root)

	*node = *root

	// 遍历找到leaf 层级二分查找
	for !node.IsLeaf {
		// 查找
		// TODO: ? B+树结构为保存最右最大节点 ?
		idx := sort.Search(len(node.Keys), func(i int) bool {
			return key <= node.Keys[i]
		})
		if idx == len(node.Keys) {
			idx = len(node.Keys) - 1
		}
		if err = t.seekNode(node, node.Children[idx]); err != nil {
			return err
		}
	}
	return nil
}

func cut(length int) int {
	return (length + 1) / 2
}

// 插入叶子节点
func insertKeyValIntoLeaf(n *Node, key uint64, rec []byte) (int, error) {
	idx := sort.Search(len(n.Keys), func(i int) bool {
		return key <= n.Keys[i]
	})
	if idx < len(n.Keys) && n.Keys[idx] == key {
		return 0, ErrorHasExistedKey
	}

	n.Keys = append(n.Keys, key)
	n.Values = append(n.Values, rec)
	// 排序 插入到合适的位置
	// TODO: 也可以用二分查找
	for i := len(n.Keys) - 1; i > idx; i-- {
		n.Keys[i] = n.Keys[i-1]
		n.Values[i] = n.Values[i-1]
	}
	n.Keys[idx] = key
	n.Values[idx] = rec
	return idx, nil
}

// 插入非叶子节点
func insertKeyValIntoNode(n *Node, key uint64, child OFFTYPE) (int, error) {
	idx := sort.Search(len(n.Keys), func(i int) bool {
		return key <= n.Keys[i]
	})
	if idx < len(n.Keys) && n.Keys[idx] == key {
		return 0, ErrorHasExistedKey
	}

	n.Keys = append(n.Keys, key)
	n.Children = append(n.Children, child)
	for i := len(n.Keys) - 1; i > idx; i-- {
		n.Keys[i] = n.Keys[i-1]
		n.Children[i] = n.Children[i-1]
	}
	n.Keys[idx] = key
	n.Children[idx] = child
	return idx, nil
}

// 移除叶子节点的数据
func removeKeyFromLeaf(leaf *Node, idx int) {
	tmpKeys := append([]uint64{}, leaf.Keys[idx+1:]...)
	leaf.Keys = append(leaf.Keys[:idx], tmpKeys...)

	tmpRecords := append([][]byte{}, leaf.Values[idx+1:]...)
	leaf.Values = append(leaf.Values[:idx], tmpRecords...)
}

// 移除非叶子节点
func removeKeyFromNode(node *Node, idx int) {
	tmpKeys := append([]uint64{}, node.Keys[idx+1:]...)
	node.Keys = append(node.Keys[:idx], tmpKeys...)

	tmpChildren := append([]OFFTYPE{}, node.Children[idx+1:]...)
	node.Children = append(node.Children[:idx], tmpChildren...)
}

// 拆分成两个叶子节点
func (t *Tree) splitLeafIntoTowleaves(leaf *Node, newLeaf *Node) error {
	var (
		i, split int
	)
	split = cut(order)

	// newLeaf为新建的叶子节点
	// insert的时候不会涉及到合并，删除的时候才会
	for i = split; i <= order; i++ {
		newLeaf.Keys = append(newLeaf.Keys, leaf.Keys[i])
		newLeaf.Values = append(newLeaf.Values, leaf.Values[i])
	}

	// adjust relation
	leaf.Keys = leaf.Keys[:split]
	leaf.Values = leaf.Values[:split]

	// 修改链接上下游
	newLeaf.Next = leaf.Next
	leaf.Next = newLeaf.Self
	newLeaf.Prev = leaf.Self
	// 共同的父节点
	newLeaf.Parent = leaf.Parent

	if newLeaf.Next != InvalidOffset {
		var (
			nextNode *Node
			err      error
		)
		// TODO: ？之前没存储么？？
		if nextNode, err = t.newMappingNodeFromPool(newLeaf.Next); err != nil {
			return err
		}
		nextNode.Prev = newLeaf.Self
		if err = t.flushNodesAndPutNodesPool(nextNode); err != nil {
			return err
		}
	}

	return err
}

func (t *Tree) insertIntoLeaf(key uint64, rec []byte) error {
	var (
		leaf    *Node
		err     error
		idx     int
		newLeaf *Node
	)

	if leaf, err = t.newMappingNodeFromPool(InvalidOffset); err != nil {
		return err
	}

	// 找到新增节点的位置
	if err = t.findLeaf(leaf, key); err != nil {
		return err
	}

	// 写入数据
	if idx, err = insertKeyValIntoLeaf(leaf, key, rec); err != nil {
		return err
	}

	// update the last key of parent's if necessary
	// 更新父节点的key
	if err = t.mayUpdatedLastParentKey(leaf, idx); err != nil {
		return err
	}

	// insert key/val into leaf
	// 符合结点数量，则不需要拆分
	if len(leaf.Keys) <= order {
		return t.flushNodeAndPutNodePool(leaf)
	}

	// split leaf so new leaf node
	// 拆分成两个叶子节点

	if newLeaf, err = t.newNodeFromDisk(); err != nil {
		return err
	}
	newLeaf.IsLeaf = true
	if err = t.splitLeafIntoTowleaves(leaf, newLeaf); err != nil {
		return err
	}
	// 落盘
	if err = t.flushNodesAndPutNodesPool(newLeaf, leaf); err != nil {
		return err
	}

	// insert split key into parent
	// 拆分成两个子节点之后，需要有一个新的父节点 存储最大值
	return t.insertIntoParent(leaf.Parent, leaf.Self, leaf.Keys[len(leaf.Keys)-1], newLeaf.Self)
}

// 获取key对应位置
func getIndex(keys []uint64, key uint64) int {
	idx := sort.Search(len(keys), func(i int) bool {
		return key <= keys[i]
	})
	return idx
}

// 新增结点插入到父节点
func insertIntoNode(parent *Node, idx int, left_off OFFTYPE, key uint64, right_off OFFTYPE) {
	var (
		i int
	)
	// 调整位置
	parent.Keys = append(parent.Keys, key)
	for i = len(parent.Keys) - 1; i > idx; i-- {
		parent.Keys[i] = parent.Keys[i-1]
	}
	parent.Keys[idx] = key

	// 跟父节点的子节点
	if idx == len(parent.Children) {
		parent.Children = append(parent.Children, right_off)
		return
	}
	tmpChildren := append([]OFFTYPE{}, parent.Children[idx+1:]...)
	parent.Children = append(append(parent.Children[:idx+1], right_off), tmpChildren...)
}

// 插入结点拆分后父节点结点树超过order
func (t *Tree) insertIntoNodeAfterSplitting(oldNode *Node) error {
	var (
		newNode, child, nextNode *Node
		err                      error
		i, split                 int
	)

	if newNode, err = t.newNodeFromDisk(); err != nil {
		return err
	}

	split = cut(order)

	for i = split; i <= order; i++ {
		newNode.Children = append(newNode.Children, oldNode.Children[i])
		newNode.Keys = append(newNode.Keys, oldNode.Keys[i])

		// update new_node children relation
		if child, err = t.newMappingNodeFromPool(oldNode.Children[i]); err != nil {
			return err
		}
		child.Parent = newNode.Self
		if err = t.flushNodesAndPutNodesPool(child); err != nil {
			return err
		}
	}
	newNode.Parent = oldNode.Parent

	oldNode.Children = oldNode.Children[:split]
	oldNode.Keys = oldNode.Keys[:split]

	newNode.Next = oldNode.Next
	oldNode.Next = newNode.Self
	newNode.Prev = oldNode.Self

	// 更新relation
	if newNode.Next != InvalidOffset {
		if nextNode, err = t.newMappingNodeFromPool(newNode.Next); err != nil {
			return err
		}
		nextNode.Prev = newNode.Self
		if err = t.flushNodesAndPutNodesPool(nextNode); err != nil {
			return err
		}
	}

	if err = t.flushNodesAndPutNodesPool(oldNode, newNode); err != nil {
		return err
	}

	// 循环更新父节点 直至不需要更新or更新到根节点
	return t.insertIntoParent(
		oldNode.Parent, oldNode.Self, oldNode.Keys[len(oldNode.Keys)-1], newNode.Self)
}

// 更新父节点 （非叶子节点）
func (t *Tree) insertIntoParent(
	parentOff OFFTYPE, leftOff OFFTYPE, key uint64, rightOff OFFTYPE) error {
	var (
		idx    int
		parent *Node
		err    error
		left   *Node
		right  *Node
	)
	// 如果更新的是根节点
	if parentOff == OFFTYPE(InvalidOffset) {
		if left, err = t.newMappingNodeFromPool(leftOff); err != nil {
			return err
		}
		if right, err = t.newMappingNodeFromPool(rightOff); err != nil {
			return err
		}
		if err = t.newRootNode(left, right); err != nil {
			return err
		}
		return t.flushNodesAndPutNodesPool(left, right)
	}

	if parent, err = t.newMappingNodeFromPool(parentOff); err != nil {
		return err
	}

	// 获取插入位置
	idx = getIndex(parent.Keys, key)
	// 更新父节点的信息
	insertIntoNode(parent, idx, leftOff, key, rightOff)

	// 落盘
	if len(parent.Keys) <= order {
		return t.flushNodesAndPutNodesPool(parent)
	}

	// 父节点也超过order，需要拆分
	return t.insertIntoNodeAfterSplitting(parent)
}

func (t *Tree) newRootNode(left *Node, right *Node) error {
	var (
		root *Node
		err  error
	)

	if root, err = t.newNodeFromDisk(); err != nil {
		return err
	}
	root.Keys = append(root.Keys, left.Keys[len(left.Keys)-1])
	root.Keys = append(root.Keys, right.Keys[len(right.Keys)-1])
	root.Children = append(root.Children, left.Self)
	root.Children = append(root.Children, right.Self)
	left.Parent = root.Self
	right.Parent = root.Self

	t.rootOff = root.Self
	return t.flushNodeAndPutNodePool(root)
}

// Insert 插入节点
func (t *Tree) Insert(key uint64, val []byte) error {
	var (
		err  error
		node *Node
	)

	// 如果是根节点 直接插入
	if t.rootOff == InvalidOffset {
		if node, err = t.newNodeFromDisk(); err != nil {
			return err
		}
		t.rootOff = node.Self
		node.IsActive = true
		node.Keys = append(node.Keys, key)
		node.Values = append(node.Values, val)
		node.IsLeaf = true
		return t.flushNodeAndPutNodePool(node)
	}

	return t.insertIntoLeaf(key, val)
}

// Update 更新节点 逻辑比较简单 直接查找后更新落盘
func (t *Tree) Update(key uint64, val []byte) error {
	var (
		node *Node
		err  error
	)

	if t.rootOff == InvalidOffset {
		return ErrorNotFoundKey
	}

	if node, err = t.newMappingNodeFromPool(InvalidOffset); err != nil {
		return err
	}

	if err = t.findLeaf(node, key); err != nil {
		return err
	}

	for i, nkey := range node.Keys {
		if nkey == key {
			node.Values[i] = val
			return t.flushNodesAndPutNodesPool(node)
		}
	}
	return ErrorNotFoundKey
}

// 如果插入的叶子节点放到了最后，需要更新父节点的值
func (t *Tree) mayUpdatedLastParentKey(leaf *Node, idx int) error {
	// update the last key of parent's if necessary
	// 如果写入的是最后一个key，那么需要更新父节点的key
	if idx == len(leaf.Keys)-1 && leaf.Parent != InvalidOffset {
		// 最后一个key值
		key := leaf.Keys[len(leaf.Keys)-1]
		updateNodeOff := leaf.Parent
		var (
			updateNode *Node
			node       *Node
		)

		// 拉取当前位置的节点
		if node, err = t.newMappingNodeFromPool(leaf.Self); err != nil {
			return err
		}
		// leaf覆盖
		*node = *leaf
		defer t.putNodePool(node)

		for updateNodeOff != InvalidOffset && idx == len(node.Keys)-1 {
			if updateNode, err = t.newMappingNodeFromPool(updateNodeOff); err != nil {
				return err
			}
			// 找到父节点中上述子节点的index
			for i, v := range updateNode.Children {
				if v == node.Self {
					idx = i
					break
				}
			}
			// 更新父节点的idx的key值
			updateNode.Keys[idx] = key
			// 落盘
			if err = t.flushNodeAndPutNodePool(updateNode); err != nil {
				return err
			}
			// 继续向上更新
			updateNodeOff = updateNode.Parent
			*node = *updateNode
		}
	}
	return nil
}

func (t *Tree) deleteKeyFromNode(off OFFTYPE, key uint64) error {
	if off == InvalidOffset {
		return nil
	}
	var (
		node      *Node
		nextNode  *Node
		prevNode  *Node
		newRoot   *Node
		childNode *Node
		idx       int
		err       error
	)
	if node, err = t.newMappingNodeFromPool(off); err != nil {
		return err
	}
	idx = getIndex(node.Keys, key)
	fmt.Printf("delete key: %d from node: %d,keys:%v,idx:%d \n", key, off, node.Keys, idx)
	removeKeyFromNode(node, idx)

	// update the last key of parent's if necessary
	if idx == len(node.Keys) {
		if err = t.mayUpdatedLastParentKey(node, idx-1); err != nil {
			return err
		}
	}

	// if statisfied len
	if len(node.Keys) >= order/2 {
		return t.flushNodesAndPutNodesPool(node)
	}

	// 根节点
	if off == t.rootOff && len(node.Keys) == 1 {
		if newRoot, err = t.newMappingNodeFromPool(node.Children[0]); err != nil {
			return err
		}
		node.IsActive = false
		newRoot.Parent = InvalidOffset
		t.rootOff = newRoot.Self
		return t.flushNodesAndPutNodesPool(node, newRoot)
	}

	if node.Next != InvalidOffset {
		if nextNode, err = t.newMappingNodeFromPool(node.Next); err != nil {
			return err
		}
		// lease from next node
		if len(nextNode.Keys) > order/2 {
			key := nextNode.Keys[0]
			child := nextNode.Children[0]

			// update child's parent
			if childNode, err = t.newMappingNodeFromPool(child); err != nil {
				return err
			}
			childNode.Parent = node.Self

			removeKeyFromNode(nextNode, 0)
			if idx, err = insertKeyValIntoNode(node, key, child); err != nil {
				return err
			}
			// update the last key of parent's if necessy
			if err = t.mayUpdatedLastParentKey(node, idx); err != nil {
				return err
			}
			return t.flushNodesAndPutNodesPool(node, nextNode, childNode)
		}
		// merge nextNode and curNode
		// 合并节点
		if node.Prev != InvalidOffset {
			if prevNode, err = t.newMappingNodeFromPool(node.Prev); err != nil {
				return err
			}
			prevNode.Next = nextNode.Self
			nextNode.Prev = prevNode.Self
			if err = t.flushNodesAndPutNodesPool(prevNode); err != nil {
				return err
			}
		} else {
			nextNode.Prev = InvalidOffset
		}

		nextNode.Keys = append(node.Keys, nextNode.Keys...)
		nextNode.Children = append(node.Children, nextNode.Children...)

		// update child's parent
		for _, v := range node.Children {
			if childNode, err = t.newMappingNodeFromPool(v); err != nil {
				return err
			}
			childNode.Parent = nextNode.Self
			if err = t.flushNodesAndPutNodesPool(childNode); err != nil {
				return err
			}
		}

		node.IsActive = false
		t.putFreeBlocks(node.Self)

		if err = t.flushNodesAndPutNodesPool(node, nextNode); err != nil {
			return err
		}

		// delete parent's key recursively
		// 向上递归更新 删除父节点的key
		return t.deleteKeyFromNode(node.Parent, node.Keys[len(node.Keys)-1])
	}

	// come here because node.Next = INVALID_OFFSET
	if node.Prev != InvalidOffset {
		if prevNode, err = t.newMappingNodeFromPool(node.Prev); err != nil {
			return err
		}
		// lease from prev leaf
		// 借一个
		if len(prevNode.Keys) > order/2 {
			key := prevNode.Keys[len(prevNode.Keys)-1]
			child := prevNode.Children[len(prevNode.Children)-1]

			// update child's parent
			if childNode, err = t.newMappingNodeFromPool(child); err != nil {
				return err
			}
			childNode.Parent = node.Self

			removeKeyFromNode(prevNode, len(prevNode.Keys)-1)
			if idx, err = insertKeyValIntoNode(node, key, child); err != nil {
				return err
			}
			// update the last key of parent's if necessy
			if err = t.mayUpdatedLastParentKey(prevNode, len(prevNode.Keys)-1); err != nil {
				return err
			}
			return t.flushNodesAndPutNodesPool(prevNode, node, childNode)
		}
		// 不够的话 和前置节点合并
		// merge prevNode and curNode
		prevNode.Next = InvalidOffset
		prevNode.Keys = append(prevNode.Keys, node.Keys...)
		prevNode.Children = append(prevNode.Children, node.Children...)

		// update child's parent
		for _, v := range node.Children {
			if childNode, err = t.newMappingNodeFromPool(v); err != nil {
				return err
			}
			childNode.Parent = prevNode.Self
			if err = t.flushNodesAndPutNodesPool(childNode); err != nil {
				return err
			}
		}

		node.IsActive = false
		t.putFreeBlocks(node.Self)

		if err = t.flushNodesAndPutNodesPool(node, prevNode); err != nil {
			return err
		}

		return t.deleteKeyFromNode(node.Parent, node.Keys[len(node.Keys)-1])
	}
	return nil
}

func (t *Tree) deleteKeyFromLeaf(key uint64) error {
	var (
		leaf     *Node
		prevLeaf *Node
		nextLeaf *Node
		err      error
		idx      int
	)
	// 获取node
	if leaf, err = t.newMappingNodeFromPool(InvalidOffset); err != nil {
		return err
	}
	// 找到key对应的叶子节点
	if err = t.findLeaf(leaf, key); err != nil {
		return err
	}

	// 获取key所在节点的index位置
	idx = getIndex(leaf.Keys, key)
	if idx == len(leaf.Keys) || leaf.Keys[idx] != key {
		t.putNodePool(leaf)
		return fmt.Errorf("not found key:%d", key)
	}

	removeKeyFromLeaf(leaf, idx)

	// if leaf is root 直接落盘更新
	if leaf.Self == t.rootOff {
		return t.flushNodesAndPutNodesPool(leaf)
	}

	// update the last key of parent's if necessary
	// 此时len已经减1，表示位置是最后一个
	if idx == len(leaf.Keys) {
		// 更新父节点的值
		if err = t.mayUpdatedLastParentKey(leaf, idx-1); err != nil {
			return err
		}
	}

	// if satisfied len
	// 当前叶子节点的长度大于order/2，满足b+树结构，不需要调整其他结构，可以直接落盘更新
	if len(leaf.Keys) >= order/2 {
		return t.flushNodesAndPutNodesPool(leaf)
	}

	// 叶子节点有后继节点
	if leaf.Next != InvalidOffset {
		if nextLeaf, err = t.newMappingNodeFromPool(leaf.Next); err != nil {
			return err
		}
		// lease from next leaf
		// 借一个节点过来
		if len(nextLeaf.Keys) > order/2 {
			key := nextLeaf.Keys[0]
			rec := nextLeaf.Values[0]
			// 删除nextLeaf的第一个key
			removeKeyFromLeaf(nextLeaf, 0)
			// 将key插入到leaf中
			if idx, err = insertKeyValIntoLeaf(leaf, key, rec); err != nil {
				return err
			}
			// update the last key of parent's if necessy
			if err = t.mayUpdatedLastParentKey(leaf, idx); err != nil {
				return err
			}
			// 落盘更新
			return t.flushNodesAndPutNodesPool(nextLeaf, leaf)
		}

		// merge nextLeaf and curleaf
		// 合并后继叶子节点
		if leaf.Prev != InvalidOffset {
			if prevLeaf, err = t.newMappingNodeFromPool(leaf.Prev); err != nil {
				return err
			}
			prevLeaf.Next = nextLeaf.Self
			nextLeaf.Prev = prevLeaf.Self
			if err = t.flushNodesAndPutNodesPool(prevLeaf); err != nil {
				return err
			}
		} else {
			nextLeaf.Prev = InvalidOffset
		}

		nextLeaf.Keys = append(leaf.Keys, nextLeaf.Keys...)
		nextLeaf.Values = append(leaf.Values, nextLeaf.Values...)
		// 释放block
		leaf.IsActive = false
		t.putFreeBlocks(leaf.Self)
		// 更新
		if err = t.flushNodesAndPutNodesPool(leaf, nextLeaf); err != nil {
			return err
		}
		// 删除父节点中已合并的key （此时key已经是更新后的key，不是待删除的旧key）
		return t.deleteKeyFromNode(leaf.Parent, leaf.Keys[len(leaf.Keys)-1])
	}

	// come here because leaf.Next = INVALID_OFFSET
	// 叶子节点没有后继节点 && 前置节点不为空
	if leaf.Prev != InvalidOffset {
		if prevLeaf, err = t.newMappingNodeFromPool(leaf.Prev); err != nil {
			return err
		}
		// lease from prev leaf
		// 从前面借一个
		if len(prevLeaf.Keys) > order/2 {
			key := prevLeaf.Keys[len(prevLeaf.Keys)-1]
			rec := prevLeaf.Values[len(prevLeaf.Values)-1]
			removeKeyFromLeaf(prevLeaf, len(prevLeaf.Keys)-1)
			if idx, err = insertKeyValIntoLeaf(leaf, key, rec); err != nil {
				return err
			}
			// update the last key of parent's if necessy
			if err = t.mayUpdatedLastParentKey(prevLeaf, len(prevLeaf.Keys)-1); err != nil {
				return err
			}
			return t.flushNodesAndPutNodesPool(prevLeaf, leaf)
		}
		// merge prevleaf and curleaf
		// 否则合并前面的叶子节点
		prevLeaf.Next = InvalidOffset
		prevLeaf.Keys = append(prevLeaf.Keys, leaf.Keys...)
		prevLeaf.Values = append(prevLeaf.Values, leaf.Values...)
		// 释放block
		leaf.IsActive = false
		t.putFreeBlocks(leaf.Self)
		// 更新
		if err = t.flushNodesAndPutNodesPool(leaf, prevLeaf); err != nil {
			return err
		}

		fmt.Println("pre keys:", prevLeaf.Keys)

		preOff := prevLeaf.Self

		// 删除父节点中已合并的key
		t.deleteKeyFromNode(leaf.Parent, leaf.Keys[len(leaf.Keys)-1])
		t.ScanTreePrint()
		var pre *Node
		if pre, err = t.newMappingNodeFromPool(preOff); err != nil {
			return err
		}

		fmt.Println("pre keys:", pre.Keys)
		// 更新前置节点的父节点的值
		t.mayUpdatedLastParentKey(pre, len(pre.Keys)-1)

		// return t.deleteKeyFromNode(leaf.Parent, prevLeaf.Keys[len(prevLeaf.Keys)-2])
	}

	return nil
}

// Delete delete the key from tree
func (t *Tree) Delete(key uint64) error {
	if t.rootOff == InvalidOffset {
		return fmt.Errorf("not found key:%d", key)
	}
	return t.deleteKeyFromLeaf(key)
}

// ScanTreePrint 遍历打印
func (t *Tree) ScanTreePrint() error {
	if t.rootOff == InvalidOffset {
		return fmt.Errorf("root = nil")
	}
	Q := make([]OFFTYPE, 0)
	Q = append(Q, t.rootOff)

	floor := 0
	var (
		curNode *Node
		err     error
	)
	for 0 != len(Q) {
		floor++

		l := len(Q)
		fmt.Printf("floor %3d:", floor)
		for i := 0; i < l; i++ {
			if curNode, err = t.newMappingNodeFromPool(Q[i]); err != nil {
				return err
			}
			defer t.putNodePool(curNode)

			// print keys
			if i == l-1 {
				fmt.Printf("%d\n", curNode.Keys)
			} else {
				fmt.Printf("%d, ", curNode.Keys)
			}
			for _, v := range curNode.Children {
				Q = append(Q, v)
			}
		}
		Q = Q[l:]
	}
	fmt.Println(strings.Repeat("-", 50))
	return nil
}
