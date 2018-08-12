package merkle

import (
	"bytes"
	"errors"
	"hash"
)

type Node struct {
	Left  *Node
	Right *Node
	Hash  []byte
	Raw   []byte
}

type Tree struct {
	Hasher hash.Hash
	Levels [][]*Node // We don't need to know a Node's parent since we're capturing all the nodes in their respective levels.
	Leaves []*Node
	Blocks [][]byte
}

func New(hasher hash.Hash, list [][]byte) (*Tree, error) {
	if len(list) == 0 {
		list = [][]byte{}
	}
	tree := &Tree{
		Hasher: hasher,
	}
	tree.AppendBlocks(list)
	return tree, nil
}

func generateLevels(tree *Tree, nodes []*Node, level int) {
	if len(nodes) == 1 {
		return
	}

	length := (len(nodes) + len(nodes)%2) / 2
	list := make([]*Node, length)

	for i := 0; i < length; i++ {
		var right *Node
		left := nodes[i*2]
		if i*2+1 >= len(nodes) {
			right = left
		} else {
			right = nodes[i*2+1]
		}

		node := &Node{
			Left:  left,
			Right: right,
			Hash:  hash_(tree.Hasher, append(left.Hash, right.Hash...)),
		}

		list[i] = node
		tree.Levels[level] = list
	}

	generateLevels(tree, list, level-1)
}

func getHeight(nodeCount int) int {
	if isPowerOf2(nodeCount) {
		return nodeCount/2 + 1
	}
	return log2(nextPowerOf2(nodeCount)) - 1
}

func hash_(hasher hash.Hash, b []byte) []byte {
	// TODO: Check for error.
	hasher.Write(b)
	defer hasher.Reset()
	return hasher.Sum(nil)
}

func isPowerOf2(n int) bool {
	return n > 0 && (n&(n-1) == 0)
}

func log2(n int) int {
	i, c := n, 0
	for i > 0 {
		i >>= 1
		c++
	}
	return c
}

func nextPowerOf2(n int) int {
	i := n
	i--
	i |= n >> 1
	i |= n >> 2
	i |= n >> 4
	i |= n >> 8
	i |= n >> 16
	i |= n >> 32
	i++
	return i
}

func (tree *Tree) AppendBlocks(blocks [][]byte) {
	nodes := make([]*Node, len(blocks))
	for i, block := range blocks {
		nodes[i] = &Node{
			Hash: hash_(tree.Hasher, block),
			Raw:  block,
		}
	}
	if len(nodes)%2 == 1 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}
	tree.Blocks = append(tree.Blocks, blocks...)
	tree.Leaves = append(tree.Leaves, nodes...)
}

func (tree *Tree) GenerateTree() error {
	if len(tree.Blocks) == 0 {
		return errors.New("[ERROR] Cannot generate tree, there are no blocks!")
	}

	nodes := tree.Leaves
	height := getHeight(len(tree.Blocks))
	tree.Levels = make([][]*Node, height)
	tree.Levels[height-1] = nodes // Leaf nodes.

	generateLevels(tree, nodes, height-1)
	return nil
}

func (tree *Tree) GetRoot() *Node {
	if len(tree.Levels) == 0 {
		return nil
	}
	return tree.Levels[0][0]
}

func (tree *Tree) IsInTree(node *Node, nextLevel int, index int) bool {
	if nextLevel == 0 {
		return true
	}

	parentIndex := index >> 1 // Determine the index of the parent node.
	side := index % 2         // Determine the position of the lookup node.

	var b []byte
	parent := tree.Levels[nextLevel][parentIndex]
	if side == 1 {
		b = append(parent.Left.Hash, node.Hash...)
	} else {
		b = append(node.Hash, parent.Right.Hash...)
	}

	if !bytes.Equal(parent.Hash, hash_(tree.Hasher, b)) {
		return false
	}

	return tree.IsInTree(parent, nextLevel-1, parentIndex)
}

// Go down the tree.
func (tree *Tree) VerifyNode(node *Node) bool {
	if node.Left == nil && node.Right == nil {
		return true
	}

	hashed := hash_(tree.Hasher, append(node.Left.Hash, node.Right.Hash...))
	if !bytes.Equal(node.Hash, hashed) {
		return false
	}

	if node.Right != nil {
		return tree.VerifyNode(node.Right)
	}

	if node.Left != nil {
		return tree.VerifyNode(node.Left)
	}

	return false
}

func (tree *Tree) VerifyProof(i interface{}) bool {
	index := -1
	switch v := i.(type) {
	case []byte:
		for idx, block := range tree.Blocks {
			if bytes.Equal(block, v) {
				index = idx
			}
		}
	case int:
		index = v
	}

	if index > -1 {
		return tree.IsInTree(tree.Leaves[index], len(tree.Levels)-1, index)
	}
	return false
}

func (tree *Tree) VerifyTree() bool {
	return tree.VerifyNode(tree.GetRoot())
}
