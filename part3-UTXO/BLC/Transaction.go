package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
)

type Transaction struct {
	TxHash []byte      //1. 交易hash
	Vins   []*TXInput  //2. 输入
	Vouts  []*TXOutput //3. 输出
}

//1. 产生创世区块时的Transaction
func NewCoinbaseTransacion(address string) *Transaction  {
	//创建创世区块交易的Vin
	txInput := &TXInput{[]byte{}, -1, "Genesis DATA"}
	//创建创世区块交易的Vout
	txOutput := &TXOutput{10, address}
	//生产交易Transaction
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	//设置Transaction的TxHash
	txCoinbase.HashTransaction()

	return txCoinbase

}



//2. 创建普通交易产生的Transaction



//将Transaction 序列化再进行 hash
func (tx *Transaction) HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())
	fmt.Printf("transationHash: %x", hash)
	tx.TxHash = hash[:]
}

func (tx *Transaction)String() string {
	var vinStrings [][]byte
	for _, vin := range tx.Vins {
		vinString := fmt.Sprint(vin)
		vinStrings = append(vinStrings, []byte(vinString))
	}
	vinString := bytes.Join(vinStrings, []byte{})

	var outStrings [][]byte
	for _, out := range tx.Vouts {
		outString := fmt.Sprint(out)
		outStrings = append(outStrings, []byte(outString))
	}

	outString := bytes.Join(outStrings, []byte{})

	return fmt.Sprintf("\tTxHash: %x, \n\t\tVins: %v, \n\t\tVout: %v", tx.TxHash, string(vinString), string(outString))
}