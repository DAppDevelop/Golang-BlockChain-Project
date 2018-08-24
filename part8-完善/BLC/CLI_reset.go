package BLC

import "os"

func (cli *CLI)Reset(nodeID string)  {
	blockchain := BlockchainObject(nodeID)
	//defer blockchain.DB.Close()

	if blockchain == nil {
		os.Exit(1)
	}

	utxoSet := &UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()
}
