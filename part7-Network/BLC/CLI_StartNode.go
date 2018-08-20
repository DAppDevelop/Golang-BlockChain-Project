package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI)startNode (nodeID string, miner string)  {
	//fmt.Println(miner)
	if miner == "" || IsValidAddress([]byte(miner)) {
		startServer(nodeID, miner)
	} else {
		fmt.Println("Miner地址无效")
		os.Exit(1)
	}
}
