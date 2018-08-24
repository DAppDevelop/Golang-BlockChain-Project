package BLC

import "os"

func (cli *CLI) printchain(nodeID string) {
	blockchain := BlockchainObject(nodeID)
	//defer blockchain.DB.Close()

	if blockchain == nil {
		os.Exit(1)
	}

	blockchain.Printchain()
}
