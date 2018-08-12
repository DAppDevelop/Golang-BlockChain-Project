package BLC

import (
	"fmt"
	"bytes"
)

type TXInput struct {
	TxID    []byte // 1. 交易的Hash
	Vout      int    //2. 存储TXOutput在Vout里面的索引(第几个交易)
	//ScriptSig string // 3. 用户名花费的是谁的钱(解锁脚本,包含数字签名)
	Signature []byte //数字签名
	PublicKey[]byte //原始公钥，钱包里的公钥
}



//判断TXInput是否指定的address消费
func (txInput *TXInput) UnlockWithAddress(pubKeyHash []byte) bool {
	pubKeyHash2:=PubKeyHash(txInput.PublicKey)
	return bytes.Compare(pubKeyHash,pubKeyHash2) == 0
}

//格式化输出
func (tx *TXInput) String() string {
	return fmt.Sprintf("\n\t\t\tTxInput_TXID: %x, Vout: %v, Signature: %x, PublicKey:%x", tx.TxID, tx.Vout, tx.Signature, tx.PublicKey)
}