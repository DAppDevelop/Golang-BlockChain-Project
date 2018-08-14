package main

import "go-BlockChain/part5-UTXOSet/BLC"

func main() {
	////创建命令行对象
	cli := BLC.CLI{}
	cli.Run()
}
