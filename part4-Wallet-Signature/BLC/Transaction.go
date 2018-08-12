package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"encoding/hex"
)

type Transaction struct {
	TxID  []byte      //1. 交易hash
	Vins  []*TXInput  //2. 输入
	Vouts []*TXOutput //3. 输出
}

//1. 产生创世区块时的Transaction
func NewCoinbaseTransacion(address string) *Transaction {
	//创建创世区块交易的Vin
	txInput := &TXInput{[]byte{}, -1, nil, nil}
	//创建创世区块交易的Vout
	//txOutput := &TXOutput{10, address}
	txOutput := NewTxOutput(10,address)
	//生产交易Transaction
	txCoinBaseTransaction := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	//设置Transaction的TxHash
	txCoinBaseTransaction.SetID()

	return txCoinBaseTransaction

}

//2. 创建普通交易产生的Transaction
func NewSimpleTransation(from string, to string, amount int64, bc *Blockchain, txs []*Transaction) *Transaction {
	//1.定义Input和Output的数组
	var txInputs []*TXInput
	var txOutputs []*TXOutput

	//获取本次转账要使用output
	total, spentableUTXO := bc.FindSpentableUTXOs(from, amount, txs)

	//获取钱包的集合：
	wallets := NewWallets()
	wallet := wallets.WalletMap[from]

	//2.创建Input
	for txID, indexArray := range spentableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TXInput{txIDBytes, index, nil, wallet.PublickKey}
			txInputs = append(txInputs, txInput)
		}
	}

	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput2 := NewTxOutput(total-amount, from)
	txOutputs = append(txOutputs, txOutput2)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	tx.SetID()
	//fmt.Println(tx)
	return tx
}

//将Transaction 序列化再进行 hash
func (tx *Transaction) SetID() {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())
	//fmt.Printf("transationHash: %x", hash)
	tx.TxID = hash[:]
}

func (tx *Transaction) IsCoinBaseTransaction() bool {
	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}

//格式化输出
func (tx *Transaction) String() string {
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

	return fmt.Sprintf("\n\r\t\t===============================\n\r\t\tTxID: %x, \n\t\tVins: %v, \n\t\tVout: %v\n\t\t", tx.TxID, string(vinString), string(outString))
}
