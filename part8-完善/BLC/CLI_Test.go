package BLC

func (cli *CLI)Test(nodeID string)  {
	bc := BlockchainObject(nodeID)
	defer bc.DB.Close()

	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()


}
