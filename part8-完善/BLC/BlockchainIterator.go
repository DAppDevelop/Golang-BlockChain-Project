package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"os"
)

type BlockchainIterator struct {
	currentHash []byte   //当前hash
	DB          *bolt.DB //数据库
}

/*
	根据当前迭代器currentHash从数据库中查找对应的block,
	之后将迭代器的currentHash置为前一个区块hash.
 */
func (blockchainIterator *BlockchainIterator) Next() *Block {
	var block Block
	DBName := fmt.Sprintf(DBName, os.Getenv("NODE_ID"))
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {
			//获取当期迭代器对应的block
			currentBlockBytes := b.Get(blockchainIterator.currentHash)
			//block = DeserializeBlock(currentBlockBytes)
			gobDecode(currentBlockBytes, &block)

			//将迭代器的currentHash 置为 上一个区块的hash
			blockchainIterator.currentHash = block.PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &block
}
