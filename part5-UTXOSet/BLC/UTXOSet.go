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

			for txIDStr, outs := range utxoMap {
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

/*
	查询对应地址的余额
 */
func (utxoSet *UTXOSet) GetBalance(address string) int64 {
	var total int64

	utxos := utxoSet.FindUnspentUTXOsByAddress(address)

	for _, utxo := range utxos{
		total += utxo.Output.Value
	}

	return total
}

/*
	查询对应地址的UTXO
 */
func (utxoSet *UTXOSet) FindUnspentUTXOsByAddress(address string) []*UTXO {
	var utxos []*UTXO

	//读数据库
	err := utxoSet.blockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOSetBucketName))
		if b != nil {
			//遍历UTXOSetBucketName 表
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				//反序列
				txOutputs := DeserializeTxOutputs(v)
				//遍历utxos
				for _, utxo := range txOutputs.UTXOs{
					//判断地址是否对应
					if utxo.Output.UnlockWithAddress(address) {
						utxos = append(utxos, utxo)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return utxos
}