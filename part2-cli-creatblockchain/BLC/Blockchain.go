package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"math/big"
)

type Blockchain struct {
	Tip []byte //最新的区块的Hash
	DB  *bolt.DB
}

// 数据库名字
const dbName = "blockchain.db"

// 表的名字
const blockTableName = "blocks"

//1. 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(data string) {

	//判断数据库是否已经存
	if DBExists() {
		fmt.Println("Genesis Block 已经存在...")
		os.Exit(1)
	}

	fmt.Println("创建创世区块....")

	//创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		//创建表
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			genesisBlock := CreateGenesisBlock(data)

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
}

// 增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(data string) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			//取到最新区块
			blockbyte := b.Get(blc.Tip)

			block := DeserializeBlock(blockbyte)

			// 创建新区块
			newBlock := NewBlock(data, block.Height + 1 , block.Hash)

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
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

// 返回Blockchain对象
func BlockchainObject() *Blockchain  {
	//因为已经知道数据库的名字，所以只要取出最新区块hash，既可以返回blockchain对象
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}


	var tip []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//取出最新区块hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	return &Blockchain{tip, db}
}
