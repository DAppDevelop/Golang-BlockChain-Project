package BLC

import (
	"fmt"
	"flag"
	"os"
	"log"
)

type CLI struct{}

func (cli *CLI) Run() {

	isValidArgs()

	//配置./moac xxx 中xxx的命令参数
	//e.g. ./moac addblock
	createblockchainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printchainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	//关联命令参数
	//sendBlockCmd
	flagFrom := sendBlockCmd.String("from", "", "转账源地址")
	flagTo := sendBlockCmd.String("to", "", "转账目的地址")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额")

	//createblockchainCmd 创世区块地址
	flagCoinbase := createblockchainCmd.String("address", "", "创世区块数据的地址")

	//getbalanceCmd
	flagGetbalanceWithAddress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.......")

	switch os.Args[1] {
	case "send":
		//解析参数
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "create":
		err := createblockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	//Parsed() -》是否执行过Parse()
	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)
	}

	if createblockchainCmd.Parsed() {
		if *flagCoinbase == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCoinbase)
	}

	if printchainCmd.Parsed() {
		cli.printchain()
	}

	if getbalanceCmd.Parsed() {
		if *flagGetbalanceWithAddress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*flagGetbalanceWithAddress)
	}


}

//输出使用指南
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreate -address --创世区块交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT --交易明细")
	fmt.Println("\tprint --输出区块信息.")
	fmt.Println("\tgetbalance -address --获取address有多少币.")
}

func (cli *CLI) createGenesisBlockchain(address string) {
	CreateBlockchainWithGenesisBlock(address)
}

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

func (cli *CLI) printchain() {
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.Printchain()
}

func (cli *CLI) getBalance (address string)  {

	fmt.Println("地址：" + address)
	//txs := UnSpentTransationsWithAdress(address)


}

//判断参数是否有效
func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
