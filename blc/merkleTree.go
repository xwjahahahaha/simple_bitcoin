package blc

import "crypto/sha256"

// merkle树（仅仅保存根结点）
type MerkleTree struct {
	RootNode *MerkleNode		// 根节点
}

// merkle节点
type MerkleNode struct {
	Left *MerkleNode
	Right *MerkleNode
	Data []byte
}

/**
 * @Description: 创建一个新的merkle节点
 * @param left
 * @param right
 * @param data
 * @return *MerkleNode
 */
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}
	// 计算当前节点数据域
	if left == nil && right == nil {		// 叶子节点
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	}else {									// 上层节点
		preHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(preHash)
		mNode.Data = hash[:]
	}
	mNode.Left = left
	mNode.Right = right
	return &mNode
}

/**
 * @Description: 创建merkle tree
 * @param data	所有的交易
 * @return *MerkleTree
 */
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	// 确保叶子节点是偶数即完美二叉树，不然就复制交易
	if len(data) % 2 != 0 {
		data = append(data, data[len(data)-1])
	}

	// 生成交易节点/叶子节点
	for _, txData := range data {
		node := NewMerkleNode(nil, nil, txData)
		nodes = append(nodes, *node)
	}

	// 依层建立节点，总层数=len(data)/2 (除去底层的交易节点/叶子节点)
	for i:=0; i<len(data)/2; i++ {
		var newLevel []MerkleNode
		// 遍历节点
		for j:=0; j<len(nodes); j+=2 {
			newNode := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *newNode)
		}
		// 更新新一层的节点
		nodes = newLevel
	}
	// 层构建完毕，最后一层的第一个节点就是merkle根
	return &MerkleTree{RootNode: &nodes[0]}
}