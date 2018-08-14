package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) send(from []string, to []string, amount []string) {
	/*
	address:  1Rs9zcPDqosXucdJjGP4wjGrtA1SmYpwGnQBMCprE2TdvhUyhk	c
	address:  1YfMAGkzTU3P19DobiAiggGzzcymvJyePughP37efhVgCV4W8e	b
	address:  1Z4DNkwSLgQR8yhtTnZSyobdenW3FfjMtkAnHJdM9ZAVenYDsU  	a
	 */
	//go run main.go send -from '["1Z4DNkwSLgQR8yhtTnZSyobdenW3FfjMtkAnHJdM9ZAVenYDsU"]' -to '["1YfMAGkzTU3P19DobiAiggGzzcymvJyePughP37efhVgCV4W8e"]' -amount '["4"]'
	//go run main.go send -from '["1Z4DNkwSLgQR8yhtTnZSyobdenW3FfjMtkAnHJdM9ZAVenYDsU","1Z4DNkwSLgQR8yhtTnZSyobdenW3FfjMtkAnHJdM9ZAVenYDsU"]' -to '["1YfMAGkzTU3P19DobiAiggGzzcymvJyePughP37efhVgCV4W8e","1Rs9zcPDqosXucdJjGP4wjGrtA1SmYpwGnQBMCprE2TdvhUyhk"]' -amount '["2","1"]'
	/*
	1/	a->b 4
	2/	a->b 2  a->c 1
	3/	b->c 1  c->a 2
	 */
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	bc := BlockchainObject()
	defer bc.DB.Close()

	bc.MineNewBlock(from, to, amount)
}