package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"os"
	"bytes"
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

			//for txid, outputs := range utxoMap {
			//	fmt.Printf("txID :%s", txid)
			//	for _, utxo := range outputs.UTXOs {
			//		fmt.Println(utxo)
			//	}
			//}

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

/*
	更新数据库UTXO
 */
func (utxoSet *UTXOSet) Update() {
	//对最后一个区块进行处理
	lastBlock := utxoSet.blockChain.Iterator().Next()

	//遍历TXs 获取所有input
	txInputs := []*TXInput{}
	for _, tx := range lastBlock.Txs {
		if !tx.IsCoinBaseTransaction() {
			for _, input := range tx.Vins {
				txInputs = append(txInputs, input)
			}
		}
	}

	//遍历TXs 获取UTXO
	outsMap := make(map[string]*TxOutputs)
	for _, tx := range lastBlock.Txs {
		//每个交易中的utxo数组
		utxos := []*UTXO{}
		for outIndex, txOut := range tx.Vouts {
			isSpent := false
			for _, txInput := range txInputs {
				if txInput.Vout == outIndex &&
					bytes.Compare(txInput.TxID, tx.TxID) == 0 {
					//已花费
					isSpent = true
					break
				}
			}
			if isSpent == false {
				utxo := &UTXO{tx.TxID, outIndex, txOut}
				utxos = append(utxos, utxo)
			}
		}

		if len(utxos) > 0 {
			txIDStr := hex.EncodeToString(tx.TxID)
			outputs := &TxOutputs{utxos}
			outsMap[txIDStr] = outputs
		}
	}

	//for txid, outputs := range outsMap {
	//	fmt.Printf("---------txID :%s", txid)
	//	for _, utxo := range outputs.UTXOs {
	//		fmt.Println(utxo)
	//	}
	//}

	//获取utxo表,将input对应的utxo删除, 添加outsMap中的utxo
	err := utxoSet.blockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOSetBucketName))
		if b != nil {
			//1.删除inputs对应的utxo
			for _, input := range txInputs {
				txOutputsBytes := b.Get(input.TxID)
				if len(txOutputsBytes) == 0 {
					continue
				}

				txOutputs := DeserializeTxOutputs(txOutputsBytes)

				//是否需要被删除
				isNeedDelete := false

				//将当前outputs里面未被消费的utxo 保存起来
				utxos := []*UTXO{}

				for _, utxo := range txOutputs.UTXOs {
					if bytes.Compare(utxo.TxID, input.TxID) == 0 &&
						input.Vout == utxo.Index &&
						input.UnlockWithAddress(utxo.Output.PubKeyHash) {
						//已花费
						isNeedDelete = true
						continue
					}

					utxos = append(utxos, utxo)
				}

				if isNeedDelete {
					err := b.Delete(input.TxID)
					if err != nil {
						log.Panic(err)
					}

					if len(utxos) > 0 {
						outputs := &TxOutputs{utxos}
						b.Put(input.TxID, outputs.Serialize())
					}
				}

			}

			//2.添加outsMap 到数据库中
			for txID, outputs := range outsMap {
				txIDBytes, _ := hex.DecodeString(txID)
				b.Put(txIDBytes, outputs.Serialize())
			}

		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
