package BLC

import "fmt"

func (cli *CLI) GetAddressList(nodeID string) {
	fmt.Println("打印所有的钱包地址。。")
	wallets := NewWallets(nodeID)
	for address, _ := range wallets.WalletMap {
		fmt.Println("address: ", address)
	}
}
