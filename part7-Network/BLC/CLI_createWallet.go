package BLC

import "fmt"

func (cli *CLI) CreateWallet(nodeID string) {
	wallets := NewWallets(nodeID)
	wallets.CreateWallet(nodeID)
	fmt.Println("钱包：", wallets.WalletMap)
}
