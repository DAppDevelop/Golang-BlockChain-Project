package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) createGenesisBlockchain(address string) {
	CreateBlockchainWithGenesisBlock(address)

	block := BlockchainObject()
	defer block.DB.Close()

	if block == nil {
		fmt.Println("没有数据库。。")
		os.Exit(1)
	}

	utxoSet := &UTXOSet{block}
	utxoSet.ResetUTXOSet()
}