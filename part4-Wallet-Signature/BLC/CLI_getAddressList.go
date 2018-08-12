package BLC

import "fmt"

func (cli *CLI) GetAddressList() {
	wallets := NewWallets()
	for address, _ := range wallets.WalletMap {
		fmt.Println("address: ", address)
	}
}
