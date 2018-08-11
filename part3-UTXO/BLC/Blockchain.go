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

	//1. 通过相关算法建立Transaction数组

	//转换amount为int
	amountInt, _ := strconv.Atoi(amount[0])

	tx := NewSimpleTransation(from[0], to[0], amountInt)
	//fmt.Println(tx)

	var txs []*Transaction
	txs = append(txs, tx)

	var block *Block
	//获取最新的block
	blockchain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = DeserializeBlock(blockBytes)
		}

		return nil
	})

	//2. 根据最新的block的信息,建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	//将新区块存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {

			b.Put(block.Hash, block.Serialize())

			b.Put([]byte("l"), block.Hash)

			blockchain.Tip = block.Hash

		}
		return nil
	})
}

// 增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))

		if b != nil {
			//取到最新区块
			blockbyte := b.Get(blc.Tip)

			block := DeserializeBlock(blockbyte)

			// 创建新区块
			newBlock := NewBlock(txs, block.Height+1, block.Hash)

			//序列号block并存入数据库
			err := b.Put(newBlock.Hash, newBlock.Serialize())

			if err != nil {
				log.Panic(err)
			}

			//更新数据库最新区块hash
			err = b.Put([]byte("l"), []byte(newBlock.Hash))

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

func (blc *Blockchain) UnSpent(address string) []*TXOutput {
	/*
	1.遍历数据库，获取每个block ---> Txs
	2.遍历所有交易：
		Inputs -- 记录为已花费
		Outputs -- 每个output
	 */

	var unSpentTxOut [] *TXOutput
	spentTxOutputMap := make(map[string][]int)

	it := blc.Iterator()

	for {

		//1、获取每个block
		block := it.Next()
		//2、遍历block的Txs
		for _, tx := range block.Txs {
			//遍历每个tx：txID/Vins/Vouts
			//遍历所有TxInput
			if !tx.IsCoinBaseTransaction() { //tx不是CoinBase交易，遍历TxInput
				for _, txInput := range tx.Vins {
					if txInput.UnlockWithAddress(address) {
						key := hex.EncodeToString(txInput.TxHash)
						spentTxOutputMap[key] = append(spentTxOutputMap[key], txInput.Vout)
					}
				}
			}

			//遍历所有的TxOutput

			for _, txOutput := range tx.Vouts {
				if txOutput.UnlockWithAddress(address) {
					if len(spentTxOutputMap) != 0 {

					} else {
						//如果长度为0，整没有花费，output无需判断
						unSpentTxOut = append(unSpentTxOut, txOutput)
					}
				}
			}
		}

		//3、判断退出
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}

	}

	return nil
}

func (blc *Blockchain) GetBalance(address string) int64 {
	txOutputs := blc.UnSpent(address)
	var total int64

	for _, txOutput := range txOutputs {
		total += txOutput.Value
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
