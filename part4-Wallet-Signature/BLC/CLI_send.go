package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) send(from []string, to []string, amount []string) {
	//go run main.go send -from '["yancey"]' -to '["a"]' -amount '["10"]'
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from, to, amount)
}