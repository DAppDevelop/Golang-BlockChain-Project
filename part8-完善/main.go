package main

import (
	"go-BlockChain/part8-完善/BLC"
	"os"
	"fmt"
	"bufio"
	"log"
)

func main() {
	//创建命令行对象
	//go signalListen()
	cli := BLC.CLI{}
	cli.Run()

}

func signalListen() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter some input: ")
	stringC := make(chan string)

	go func(stringC chan string) {
		for {
			i := <-stringC
			switch i {
			case "Philip\n", "Ivo\n", "Chris\n":
				fmt.Printf("Welcome %s\n", i)
			default:
				fmt.Println("You are not welcome here! Goodbye!")
			}
		}
	}(stringC)

	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		stringC <- input

	}

}
