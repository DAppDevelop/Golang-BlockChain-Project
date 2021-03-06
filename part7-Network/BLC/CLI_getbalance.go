package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) getBalance(address string,nodeID string) {
	blockchain := BlockchainObject(nodeID)
	if blockchain == nil {
		os.Exit(1)
	}
	defer blockchain.DB.Close()
	//txs 传nil值，查询时没有新的交易产生
	//total := blockchain.GetBalance(address, []*Transaction{})
	utxoSet := &UTXOSet{blockchain}
	total := utxoSet.GetBalance(address)
	fmt.Printf("%s的余额：%d\n", address, total)
}
