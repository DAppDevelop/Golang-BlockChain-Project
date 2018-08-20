package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	currentHash []byte
	DB *bolt.DB
}

func (blockchainIterator *BlockchainIterator)Next() *Block  {
	var block Block

	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
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
