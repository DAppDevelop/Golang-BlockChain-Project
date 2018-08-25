package BLC

import (
	"fmt"
	"bytes"
)

type TXOutput struct {
	Value        int64  //金额
	//ScriptPubKey string //用户名(scriptPubkey:锁定脚本,包含公钥)
	PubKeyHash [] byte//公钥哈希
}

//判断TxOutput是否时指定的用户解锁
func (txOutput *TXOutput) UnlockWithAddress(address string) bool{
	full_payload:=Base58Decode([]byte(address))

	pubKeyHash:=full_payload[1:len(full_payload)-addressCheckSumLen]

	return bytes.Compare(pubKeyHash,txOutput.PubKeyHash) == 0
}

//根据地址创建一个output对象
func NewTxOutput(value int64,address string) *TXOutput{
	txOutput:=&TXOutput{value,nil}
	txOutput.Lock(address)
	return txOutput
}

//锁定
func (tx *TXOutput) Lock(address string){
	full_payload := Base58Decode([]byte(address))
	//获取公钥hash
	tx.PubKeyHash = full_payload[1:len(full_payload)-addressCheckSumLen]
}

//格式化输出
func (tx *TXOutput) String() string {
	return fmt.Sprintf("\n\t\t\tValue: %d, PubKeyHash(转成地址显示): %s", tx.Value, PublicHashToAddress(tx.PubKeyHash))
	//return fmt.Sprintf("\n\t\t\tValue: %d, PubKeyHash: %x", tx.Value, tx.PubKeyHash)
}
