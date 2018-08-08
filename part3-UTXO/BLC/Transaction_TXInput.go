package BLC

import "fmt"

type TXInput struct {
	TxHash    []byte // 1. 交易的Hash
	Vout      int    //2. 存储TXOutput在Vout里面的索引(第几个交易)
	ScriptSig string // 3. 用户名花费的是谁的钱(解锁脚本,包含数字签名)
}


func (tx *TXInput)String() string {
	return fmt.Sprintf("TxHash: %x, Vout: %v, ScriptSig: %v", tx.TxHash, tx.Vout, tx.ScriptSig)
}