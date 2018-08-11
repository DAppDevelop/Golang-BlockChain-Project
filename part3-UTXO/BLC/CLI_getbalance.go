package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) getBalance(address string) {
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	total := blockchain.GetBalance(address)
	fmt.Printf("%s的余额：%d", address, total)

	//txs := UnSpentTransationsWithAdress(address)

}
