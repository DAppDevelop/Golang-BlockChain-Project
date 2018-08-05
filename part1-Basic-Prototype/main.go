package main

import (
	"github.com/boltdb/bolt"
	"log"
	"blockchain/part1-Basic-Prototype/BLC"
	"fmt"
)

func main() {

	//-------------------1.创建区块链、添加区块到区块链中
	//blockchain := BLC.CreatBlockchainWithGenesisBlock()
	//
	//blockchain.AddBlockToBlockchain(
	//	"second Block",
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain(
	//	"3 Block",
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain(
	//	"4 Block",
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,
	//	blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//fmt.Println(blockchain)

	//--------------------2.判断工作量证明是否有效
	//block := BLC.NewBlock("hehhehe", 0, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,})
	//fmt.Println(block)
	//
	//proofOfWork := BLC.NewProofOfWork(block)
	//fmt.Println(proofOfWork.IsValid())

	//--------------------3.创建区块并将其序列化并保存到数据库中
	//block := BLC.NewBlock(
	//	"hehhehe",
	//	0,
	//	[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,})
	//
	//fmt.Println(block)

	//打开数据库
	db, err := bolt.Open("block.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//创建表
	//err = db.Update(func(tx *bolt.Tx) error {
	//	//直接获取表，如果不存在，创建
	//	b := tx.Bucket([]byte("blockchain"))
	//
	//	if b == nil {
	//		b, err = tx.CreateBucket([]byte("blockchain"))
	//		if err != nil {
	//			log.Panic("block table create failed...")
	//		}
	//	}
	//
	//	err := b.Put([]byte("l"), block.Serialize())
	//
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//
	//	fmt.Println("写入成功！")
	//	return nil
	//})
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	//读取表里面的block

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("blockchain"))

		if b == nil {
			log.Panic("读取不到blockchain")
		}

		blockData := b.Get([]byte("l"))

		block := BLC.DeserializeBlock(blockData)

		fmt.Println("READBLOCK ", block)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
