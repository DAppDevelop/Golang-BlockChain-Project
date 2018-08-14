package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
)

type UTXOSet struct {
	blockChain *Blockchain
}

//提供一个重置的功能：获取blockchain中所有的未花费utxo

/*
查询block块中所有的未花费utxo：执行FindUnspentUTXOMap--->map

 */

func (utxoset *UTXOSet) ResetUTXOSet() {

	err := utxoset.blockChain.DB.Update(func(tx *bolt.Tx) error {
		//1.utxoset表存在，删除
		b := tx.Bucket([]byte(UTXOSetBucketName))
		if b != nil {
			err := tx.DeleteBucket([]byte(UTXOSetBucketName))
			if err != nil {
				log.Panic(err)
			}
		}

		b, err := tx.CreateBucket([]byte(UTXOSetBucketName))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			utxoMap := utxoset.blockChain.FindUnspentUTXOMap()

			for txid, outputs := range utxoMap {
				fmt.Printf("txID :%s", txid)
				for _, utxo := range outputs.UTXOs {
					fmt.Println(utxo)
				}
			}

			for txIDStr, outs := range utxoMap{
				txID, _ := hex.DecodeString(txIDStr)
				b.Put(txID, outs.Serialize())
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
