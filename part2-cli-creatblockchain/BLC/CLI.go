package BLC

import (
	"fmt"
	"flag"
	"os"
	"log"
)

type CLI struct {}

func (cli *CLI) Run() {

	isValidArgs()

	//配置./moac xxx 中xxx的命令参数
	//e.g. ./moac addblock
	addblockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	createblockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	//关联命令参数
	flagAddBlockData := addblockCmd.String("data", "chenysh", "交易数据")
	flagCreateBlockchainWithCmd := createblockchainCmd.String("data", "GenesisBlock.......", "创世区块数据")

	switch os.Args[1] {
	case "addblock":
		//解析参数
		err := addblockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createblockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	//Parsed() -》是否执行过Parse()
	if addblockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}

		cli.addBlock(*flagAddBlockData)
	}

	if createblockchainCmd.Parsed() {
		if *flagCreateBlockchainWithCmd == "" {
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockchainWithCmd)
	}

	if printchainCmd.Parsed() {
		cli.printchain()
	}

}

//输出使用指南
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -data -- 创世区块交易数据.")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")
}


func (cli *CLI) addBlock(data string) {
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	blockchain.AddBlockToBlockchain(data)

}

func (cli *CLI) printchain() {
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	blockchain.Printchain()

}

func (cli *CLI) createGenesisBlockchain(data string) {
	CreateBlockchainWithGenesisBlock(data)
}

//判断参数是否有效
func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
