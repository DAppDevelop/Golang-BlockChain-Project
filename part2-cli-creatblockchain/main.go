package main

import "blockchain/part2-cli-creatblockchain/BLC"

func main() {
	////创建命令行对象
	cli := BLC.CLI{}
	cli.Run()
}