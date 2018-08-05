package main

import (
	"blockchain/part1-Basic-Prototype/BLC"
	"fmt"
)

func main() {
	//blockchain := BLC.CreatBlockchainWithGenesisBlock()
	//
	//blockchain.AddBlockToBlockchain(
	//	"second Block",
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain(
	//	"3 Block",
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain(
	//	"4 Block",
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//fmt.Println(blockchain)



	//判断工作量证明是否有效
	//block := BLC.NewBlock("hehhehe", 0, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,})
	//fmt.Println(block)
	//
	//proofOfWork := BLC.NewProofOfWork(block)
	//fmt.Println(proofOfWork.IsValid())

	//序列化
	block := BLC.NewBlock("hehhehe", 0, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,})

	fmt.Println(block)

	bytes := block.Serialize()

	fmt.Println(bytes)

	block = BLC.DeserializeBlock(bytes)

	fmt.Println(block)
}
