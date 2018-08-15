package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"os"
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

	for _, utxo := range utxos {
		total += utxo.Output.Value
	}

	return total
}

/*
	查询对应地址, 已打包的UTXO
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
				for _, utxo := range txOutputs.UTXOs {
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

/*
	查询本次转账将要使用的UTXO
 */
func (utxoSet *UTXOSet) FindSpentableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	var total int64
	spentableUTXOMap := make(map[string][]int)

	//未打包的UTXO
	unPacketUTXO := utxoSet.FindUnpacketUTXO(from, txs)

	for _, utxo := range unPacketUTXO {
		total += utxo.Output.Value
		txIDStr := hex.EncodeToString(utxo.TxID)
		spentableUTXOMap[txIDStr] = append(spentableUTXOMap[txIDStr], utxo.Index)

		if total >= amount {
			return total, spentableUTXOMap
		}

	}

	//已在区块的UTXO
	packedUTXO := utxoSet.FindUnspentUTXOsByAddress(from)

	for _, utxo := range packedUTXO {
		total += utxo.Output.Value
		txIDStr := hex.EncodeToString(utxo.TxID)
		spentableUTXOMap[txIDStr] = append(spentableUTXOMap[txIDStr], utxo.Index)

		if total >= amount {
			return total, spentableUTXOMap
		}
	}

	if total < amount {
		fmt.Printf("%s 的余额不足, 无法转账. 余额为: %d", from, total)
		os.Exit(1)
	}


	return total, spentableUTXOMap
}

/*
	查找对应地址,未打包的UTXO
 */
func (utxoSet *UTXOSet) FindUnpacketUTXO(from string, txs []*Transaction) []*UTXO {

	//存储未花费的TxOutput
	var utxos [] *UTXO
	//存储已经花费的信息
	spentTxOutputMap := make(map[string][]int) // map[TxID] = []int{vout}

	for i := len(txs) - 1; i >= 0; i-- {
		tx := txs[i]

		utxos = caculate(tx, from, spentTxOutputMap, utxos)
	}

	return utxos
}
