package BLC

import (
	"crypto/sha256"
)

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	DataHash  []byte
}

type MerkleTree struct {
	RootNode *MerkleNode
}

func NewMerkleNode(lelfNode , rightNode *MerkleNode, txHash []byte) *MerkleNode  {
	node := &MerkleNode{}

	var hash [32]byte
	if lelfNode == nil && rightNode == nil {
		//如果是叶子节点
		hash = sha256.Sum256(txHash)
	} else {
		//子节点
		//拼接dataHash
		prevHash := append(lelfNode.DataHash, rightNode.DataHash...)
		hash = sha256.Sum256(prevHash)
	}

	node.LeftNode = lelfNode
	node.RightNode = rightNode
	node.DataHash = hash[:]

	return node

}

func NewMerkleTree(txHashData [][]byte) *MerkleTree  {
	//保存每层merkle Tree 节点 当节点数为1时, 跳出循环
	var nodes []*MerkleNode

	//创建叶子节点
	//当txHashData 为奇数 , 最后一个复制补全
	if len(txHashData) % 2 != 0 {
		txHashData = append(txHashData, txHashData[len(txHashData)-1])
	}

	for _, txHash := range txHashData{
		node := NewMerkleNode(nil, nil, txHash)
		nodes = append(nodes, node)
	}

	//生成子节点(循环到根节点生成为止)
	for {
		//每次循环新建一个newNodes 保存此层node的数组
		var newNodes []*MerkleNode

		for i := 0 ; i < len(nodes); i += 2 {
			node := &MerkleNode{nodes[i], nodes[i+1], nil}
			newNodes = append(newNodes, node)
		}

		//设置新的nodes
		nodes = newNodes

		//判断当前层node的数量是否为1, 为1则为根节点
		if len(newNodes) == 1 {
			break
		}
	}

	merkleTree := &MerkleTree{nodes[0]}

	return merkleTree
}

//func getCircleCount(len int) int {
//	count := 0
//	for {
//		if int(math.Pow(2, float64(count))) >= len {
//			return count
//		}
//		count++
//	}
//}