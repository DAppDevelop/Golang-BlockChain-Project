package main

import (
	"go-BlockChain/part2-cli-creatblockchain/BLC"
)

func main() {
	//blockchain := BLC.CreatBlockchainWithGenesisBlock()
	//fmt.Println(blockchain)

	cli := BLC.CLI{}
	cli.Run()
}
