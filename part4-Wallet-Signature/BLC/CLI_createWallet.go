package BLC

import "fmt"

func (cli *CLI) CreateWallet() {
	wallets := NewWallets()
	wallets.CreateWallet()
	fmt.Println("钱包：", wallets.WalletMap)
}
