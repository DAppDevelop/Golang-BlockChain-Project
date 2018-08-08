package main

import "go-BlockChain/part3-UTXO/BLC"

func main() {
	////创建命令行对象
	cli := BLC.CLI{}
	cli.Run()

}