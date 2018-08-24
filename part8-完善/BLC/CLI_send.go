package BLC

import (
	"os"
	"fmt"
)

func (cli *CLI) send(from []string, to []string, amount []string, nodeID string, mine bool) {
	/*
	address:  1Rs9zcPDqosXucdJjGP4wjGrtA1SmYpwGnQBMCprE2TdvhUyhk	c
	address:  1YfMAGkzTU3P19DobiAiggGzzcymvJyePughP37efhVgCV4W8e	b
	address:  1Z4DNkwSLgQR8yhtTnZSyobdenW3FfjMtkAnHJdM9ZAVenYDsU  	a
	 */
	//go run main.go send -from '["1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa"]' -to '["1RAbXZJVTYvPfdrXRd274vjrBE6XWxeyMLSxNYZixsuT7Uetrc"]' -amount '["1"]'
	//go run main.go send -from '["1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa","1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa"]' -to '["1RAbXZJVTYvPfdrXRd274vjrBE6XWxeyMLSxNYZixsuT7Uetrc","1RAbXZJVTYvPfdrXRd274vjrBE6XWxeyMLSxNYZixsuT7Uetrc"]' -amount '["2","1"]'
	//go run main.go send -from '["1YfMAGkzTU3P19DobiAiggGzzcymvJyePughP37efhVgCV4W8e","1Rs9zcPDqosXucdJjGP4wjGrtA1SmYpwGnQBMCprE2TdvhUyhk"]' -to '["1Rs9zcPDqosXucdJjGP4wjGrtA1SmYpwGnQBMCprE2TdvhUyhk","1Z4DNkwSLgQR8yhtTnZSyobdenW3FfjMtkAnHJdM9ZAVenYDsU"]' -amount '["3","1"]'
	//go run main.go send -from '["1Z4DNkwSLgQR8yhtTnZSyobdenW3FfjMtkAnHJdM9ZAVenYDsU"]' -to '["1Rs9zcPDqosXucdJjGP4wjGrtA1SmYpwGnQBMCprE2TdvhUyhk"]' -amount '["8"]'
	/*
	1/	a->b 4					a: 16 / b: 4 / c: 0
	2/	a->b 2  a->c 1			a: 23 / b: 6 / c: 1
	3/	b->c 3  c->a 1			a: 24 / b: 13 / c: 3
	4/  a->c 8					a: 26 / b: 13 / c: 11
	 */

	blockchain := BlockchainObject(nodeID)
	//defer blockchain.DB.Close()

	if blockchain == nil {
		os.Exit(1)
	}

	if mine {
		fmt.Println("--------------本地挖矿")
		txs := blockchain.hanldeTransations(from, to, amount, nodeID)
		newBlock := blockchain.MineNewBlock(txs)
		blockchain.SaveNewBlockToBlockchain(newBlock)
		utxoSet := &UTXOSet{blockchain}
		utxoSet.Update()
	} else {
		fmt.Println("--------------挖矿节点挖矿")
		//拼接nodeID到ip后
		nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
		txs := blockchain.hanldeTransations(from, to, amount, nodeID)
		//fmt.Println(nodeAddress)
		if nodeAddress != knowNodes[0] {
			//非主节点的交易先发送给主节点
			fmt.Println("sendTransactionToMainNode")
			sendTransactionToMainNode(knowNodes[0], txs)
		} else {
			//如果交易是在主节点, 直接发送给挖矿节点
			fmt.Println("sendTransactionToMiner")
			sendTransactionToMiner(knowNodes[1], txs)
		}
	}

}

/*
	3000:
	address:  12Mqbd1ALZQ9k4kqDvuYkCJgcSgAJNdRLVz7eWx43Ddq5jScG3p
	address:  1fQVvtoGtUeBNj2VvfpVL4ouB2jQgW5LWbzLh5LCYwpX4GXgDJ
	address:  1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa
--------------------------------------------------------------------
	3001:
	address:  121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x
	address:  1RAbXZJVTYvPfdrXRd274vjrBE6XWxeyMLSxNYZixsuT7Uetrc
 */


/*
	主节点转账: 初始  1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa: 10

	go run main.go send -from '["1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa"]' -to '["121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x"]' -amount '["1"]' -mine f
	go run main.go send -from '["1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa","1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa"]' -to '["121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x","121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x"]' -amount '["2","1"]'




	go run main.go send -from '["121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x"]' -to '["1RAbXZJVTYvPfdrXRd274vjrBE6XWxeyMLSxNYZixsuT7Uetrc"]' -amount '["1"]' -mine f








*/