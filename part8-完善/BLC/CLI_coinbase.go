package BLC

import (
	"io/ioutil"
	"fmt"
	"log"
	"strings"
)

func (cli *CLI)setCoinbase(nodeID string, coinbase string)  {
	//将coinbase写入文件
	fileName := fmt.Sprintf("coinbase_%s", nodeID)
	data :=  []byte(coinbase)
	if ioutil.WriteFile(fileName,data,0644) == nil {
		fmt.Println("写入文件成功:",coinbase)
	}
}

func CoinbaseAddress(nodeID string) string {
	fileName := fmt.Sprintf("coinbase_%s", nodeID)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panic(err)
	}

	result := strings.Replace(string(b),"\n","",1)
	//fmt.Println("result :", result)
	return result
}
