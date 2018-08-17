package BLC

func (cli *CLI)Test()  {
	bc := BlockchainObject()
	defer bc.DB.Close()

	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()


}
