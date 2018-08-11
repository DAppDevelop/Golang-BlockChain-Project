package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"encoding/hex"
)

type Blockchain struct {
	Tip []byte //最新的区块的Hash
	DB  *bolt.DB
}

//1. 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string) {

	//判断数据库是否已经存
	if DBExists() {
		fmt.Println("Genesis Block 已经存在...")
		os.Exit(1)
	}

	fmt.Println("创建创世区块....")

	//创建或打开数据库
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		//创建表
		b, err := tx.CreateBucket([]byte(BlockBucketName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建了一个coinbase Transaction
			txCoinbase := NewCoinbaseTransacion(address)
			// 创建创世区块
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})

			//序列号block并存入数据库
			err := b.Put([]byte(genesisBlock.Hash), []byte(genesisBlock.Serialize()))

			if err != nil {
				log.Panic(err)
			}

			//更新数据库最新区块hash
			err = b.Put([]byte("l"), []byte(genesisBlock.Hash))

			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

// 挖矿产生区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {
	/*
	1.新建交易
	2.新建区块：
		读取数据库，获取最后一块block
	3.存入到数据库中
	 */

	//1. 通过相关算法建立Transaction数组
	var txs []*Transaction
	for i := 0; i < len(from); i++ {
		//转换amount为int
		amountInt, _ := strconv.Atoi(amount[i])
		tx := NewSimpleTransation(from[i], to[i], int64(amountInt), blockchain, txs)
		//fmt.Println(tx)
		txs = append(txs, tx)
	}

	var block *Block
	//获取最新的block
	err := blockchain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = DeserializeBlock(blockBytes)
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	//2. 根据最新的block的信息,建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	//将新区块存储到数据库
	err = blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {

			b.Put(block.Hash, block.Serialize())

			b.Put([]byte("l"), block.Hash)

			blockchain.Tip = block.Hash

		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// 增加区块到区块链里面
//func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {
//
//	err := blc.DB.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(BlockBucketName))
//
//		if b != nil {
//			//取到最新区块
//			blockbyte := b.Get(blc.Tip)
//
//			block := DeserializeBlock(blockbyte)
//
//			// 创建新区块
//			newBlock := NewBlock(txs, block.Height+1, block.Hash)
//
//			//序列号block并存入数据库
//			err := b.Put(newBlock.Hash, newBlock.Serialize())
//
//			if err != nil {
//				log.Panic(err)
//			}
//
//			//更新数据库最新区块hash
//			err = b.Put([]byte("l"), []byte(newBlock.Hash))
//
//			if err != nil {
//				log.Panic(err)
//			}
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		log.Panic(err)
//	}
//
//}

//设计一个方法，用于获取指定用户的所有的未花费Txoutput
/*
UTXO模型：未花费的交易输出
	Unspent Transaction TxOutput
 */
func (blc *Blockchain) UnSpent(address string, txs []*Transaction) []*UTXO {
	/*
	0.查询本次转账已经创建了的哪些transaction
	1.遍历数据库，获取每个block ---> Txs
	2.遍历所有交易：
		Inputs -- 记录为已花费
		Outputs -- 每个output
	 */
	//存储未花费的TxOutput
	var unSpentUTXOs [] *UTXO
	//存储已经花费的信息
	spentTxOutputMap := make(map[string][]int) // map[TxID] = []int{vout}

	//第一部分：先查询本次转账，已经产生了的Transanction
	for i := len(txs) - 1; i >= 0; i-- {
		unSpentUTXOs = caculate(txs[i], address, spentTxOutputMap, unSpentUTXOs)
	}

	it := blc.Iterator()

	for {

		//1、获取每个block
		block := it.Next()
		//2、遍历block的Txs
		//倒序遍历Transactions
		for i := len(block.Txs) - 1; i >= 0; i-- {
			unSpentUTXOs = caculate(block.Txs[i], address, spentTxOutputMap, unSpentUTXOs)
		}

		//3、判断退出
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}

	}

	return unSpentUTXOs
}

//计算对应address的未花费TXOutput
func caculate(tx *Transaction, address string, spentTxOutputMap map[string][]int, unSpentUTXOs []*UTXO) []*UTXO {
	//遍历每个tx：txID，Vins，Vouts

	//遍历所有的TxInput
	if !tx.IsCoinBaseTransaction() { //tx不是CoinBase交易，遍历TxInput
		for _, txInput := range tx.Vins {
			//txInput-->TxInput
			if txInput.UnlockWithAddress(address) {
				//txInput的解锁脚本(用户名) 如果和钥查询的余额的用户名相同，
				key := hex.EncodeToString(txInput.TxID)
				spentTxOutputMap[key] = append(spentTxOutputMap[key], txInput.Vout)
				/*
				map[key]-->value
				map[key] -->[]int
				 */
			}
		}
	}

	//遍历所有的TxOutput
outputs:
	for index, txOutput := range tx.Vouts { //index= 0,txoutput.锁定脚本：王二狗
		if txOutput.UnlockWithAddress(address) {
			if len(spentTxOutputMap) != 0 {
				var isSpentOutput bool //false
				//遍历map
				for txID, indexArray := range spentTxOutputMap { //143d,[]int{1}
					//遍历 记录已经花费的下标的数组
					for _, i := range indexArray {
						if i == index && hex.EncodeToString(tx.TxID) == txID {
							isSpentOutput = true //标记当前的txOutput是已经花费
							continue outputs
						}
					}
				}

				if !isSpentOutput {
					//unSpentTxOutput = append(unSpentTxOutput, txOutput)
					//根据未花费的output，创建utxo对象--->数组
					utxo := &UTXO{tx.TxID, index, txOutput}
					unSpentUTXOs = append(unSpentUTXOs, utxo)
				}

			} else {
				//如果map长度未0,证明还没有花费记录，output无需判断
				//unSpentTxOutput = append(unSpentTxOutput, txOutput)
				utxo := &UTXO{tx.TxID, index, txOutput}
				unSpentUTXOs = append(unSpentUTXOs, utxo)
			}
		}
	}
	return unSpentUTXOs

}

/*
提供一个方法，返回用于一次转账的交易中，即将被使用为花费的utxo
 */
func (bc *Blockchain) FindSpentableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	/*
	1.根据from获取到的所有的utxo
	2.遍历utxos，累加余额，判断，是否如果余额，大于等于要要转账的金额，


	返回：map[txID] -->[]int{下标1，下标2} --->Output
	 */
	var total int64

	spentableMap := make(map[string][]int)
	//1.获取所有的utxo ：10
	utxos := bc.UnSpent(from, txs)
	//2.找即将使用utxo：3个utxo
	for _, utxo := range utxos {
		total += utxo.Output.Value
		txIDstr := hex.EncodeToString(utxo.TxID)
		spentableMap[txIDstr] = append(spentableMap[txIDstr], utxo.Index)

		if total >= amount {
			break
		}
	}

	//3.判断total是否大于等于amount
	if total < amount {
		fmt.Printf("%s，余额不足，无法转账。。", from)
		os.Exit(1)
	}

	return total, spentableMap

}

//提供一个功能：查询余额
func (blc *Blockchain) GetBalance(address string, txs []*Transaction) int64 {
	unSpentUTXOs := blc.UnSpent(address, txs)
	var total int64

	for _, utxo := range unSpentUTXOs {
		total += utxo.Output.Value
	}

	return total
}

// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain() {
	//创建迭代器
	blockIterator := blc.Iterator()

	for {
		block := blockIterator.Next()

		fmt.Println(block)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		//判断当期的block是否为创世区块（创世区块perblockhash为000000....）
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

//迭代器
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		return false
	}

	return true
}

// 返回Blockchain对象
func BlockchainObject() *Blockchain {
	//因为已经知道数据库的名字，所以只要取出最新区块hash，既可以返回blockchain对象
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {
			//取出最新区块hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	return &Blockchain{tip, db}
}
