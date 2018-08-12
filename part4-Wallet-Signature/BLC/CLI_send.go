package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) send(from []string, to []string, amount []string) {
	//go run main.go send -from '["12Gb7Fc3PeqUQMvFJcbrKD5THTpCSUPp8DNYnbUfCxRVx5EzuQt"]' -to '["12Gi79PZ7JrxkEESb3qj3RMRjHHQCn7X63xXVB3D4CpTg5igfYy"]' -amount '["2"]'
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from, to, amount)
}