package BLC

import "fmt"

func (cli *CLI)startNode (miner string)  {
	fmt.Println(miner)
	startServer("", miner)
}
