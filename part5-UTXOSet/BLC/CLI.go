package BLC

import (
	"fmt"
	"flag"
	"os"
	"log"
)

type CLI struct{}

func (cli *CLI) Run() {

	/*
	Usage:
		addblock -data DATA
		printchain


	./bc printchain
		-->执行打印的功能

	./bc send -from '["yancey"]' -to '["alice"]' -amount '["11"]'
	./bc send -from '["yancey","alice"]' -to '["bob","cici"]' -amount '["4","5"]'

	 */

	isValidArgs()

	//1.---------创建flagset命令对象
	//e.g. ./moac addblock
	//./bc  命令 -参数名 参数
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	getAddresslistsCmd:=flag.NewFlagSet("getaddresslists",flag.ExitOnError)
	createblockchainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	//2.----------设置命令后的参数对象
	flagFrom := sendCmd.String("from", "", "转账源地址")
	flagTo := sendCmd.String("to", "", "转账目的地址")
	flagAmount := sendCmd.String("amount", "", "转账金额")

	//createblockchainCmd 创世区块地址
	flagCoinbase := createblockchainCmd.String("address", "", "创世区块数据的地址")

	//getbalanceCmd
	flagGetbalanceWithAddress := getBalanceCmd.String("address", "", "要查询余额的账户.......")

	//3.----------解析参数
	switch os.Args[1] {
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "create":
		if err := createblockchainCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "print":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "createwallet":
		if err := createWalletCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "getaddresslists":
		if err := getAddresslistsCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	//4.---------根据终端输入的命令执行对应的功能
	//Parsed() -》是否执行过Parse()
	if sendCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from, to, amount)
	}

	if createblockchainCmd.Parsed() {
		if *flagCoinbase == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCoinbase)
	}

	if printChainCmd.Parsed() {
		cli.printchain()
	}

	if getBalanceCmd.Parsed() {
		if *flagGetbalanceWithAddress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*flagGetbalanceWithAddress)
	}

	if createWalletCmd.Parsed() {
		cli.CreateWallet()
	}

	if getAddresslistsCmd.Parsed() {
		cli.GetAddressList()
	}

}

//输出使用指南
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreatewallet -- 创建钱包")
	fmt.Println("\tgetaddresslists -- 获取所有的钱包地址")
	fmt.Println("\tcreate -address --创世区块交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT --交易明细")
	fmt.Println("\tprint --输出区块信息.")
	fmt.Println("\tgetbalance -address --获取address有多少币.")
}

//判断参数是否有效
func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
