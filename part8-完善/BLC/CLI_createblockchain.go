package BLC

import (
	"os"
)

func (cli *CLI) createGenesisBlockchain(address string, nodeID string) {
	CreateBlockchainWithGenesisBlock(address,nodeID)

	blockchain := BlockchainObject(nodeID)
	defer blockchain.DB.Close()

	if blockchain == nil {
		os.Exit(1)
	}

	utxoSet := &UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()
}