package BLC

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	DataHash  []byte
}

type MerkleTree struct {
	RootNode *MerkleNode
}

