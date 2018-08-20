package BLC

import "os"

func (cli *CLI) printchain(nodeID string) {
	blockchain := BlockchainObject(nodeID)
	if blockchain == nil{
		os.Exit(1)
	}
	defer blockchain.DB.Close()

	blockchain.Printchain()
}
