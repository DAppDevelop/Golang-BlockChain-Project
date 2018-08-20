package BLC

import (
	"os"
)

func (cli *CLI) createGenesisBlockchain(address string, nodeID string) {
	CreateBlockchainWithGenesisBlock(address,nodeID)

	block := BlockchainObject(nodeID)
	if block == nil {
		os.Exit(1)
	}
	defer block.DB.Close()

	utxoSet := &UTXOSet{block}
	utxoSet.ResetUTXOSet()
}