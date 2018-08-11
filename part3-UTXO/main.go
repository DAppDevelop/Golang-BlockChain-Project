package main

import (
	"blockchain/part3-UTXO/BLC"
)

func main() {
	////创建命令行对象
	cli := BLC.CLI{}
	cli.Run()

}