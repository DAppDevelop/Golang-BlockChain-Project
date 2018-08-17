package BLC

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

type Wallets struct {
	WalletMap map[string]*Wallet
}

const walletsFile = "Wallets.dat" //存储钱包数据的本地文件名

//提供一个函数，用于创建一个钱包的集合
/*
思路：修改该方法：
	读取本地的钱包文件，如果文件存在，直接获取
	如果文件不存在，创建钱包对象
 */

func NewWallets() *Wallets {
	//step1：钱包文件不存在
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		fmt.Println("钱包文件不存在。。。")
		wallets := &Wallets{}
		wallets.WalletMap = make(map[string]*Wallet)
		return wallets
	}

	wsBytes, err := ioutil.ReadFile(walletsFile)
	if err != nil {
		log.Panic(err)
	}

	gob.Register(elliptic.P256())
	var wallets Wallets

	reader := bytes.NewReader(wsBytes)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&wallets)
	if err != nil {

		log.Panic(err)
	}

	return &wallets
}

func (ws *Wallets) CreateWallet()  {
	wallet := NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("创建的钱包地址：%s\n",address)

	ws.WalletMap[string(address)] =wallet

	ws.saveFile()
}

func (ws *Wallets) saveFile () {
	var buf bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(ws)

	if err != nil {
		log.Panic(err)
	}

	wsBytes := buf.Bytes()

	err = ioutil.WriteFile(walletsFile, wsBytes, 0644)
	if err != nil {
		log.Panic(err)
	}

}
