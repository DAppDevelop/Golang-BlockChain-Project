package BLC

import "fmt"

type TXOutput struct {
	Value        int64  //金额
	ScriptPubKey string //用户名(scriptPubkey:锁定脚本,包含公钥)
}

//判断TXOutput是否指定的address解锁
func (txOutput *TXOutput) UnlockWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}

//格式化输出
func (tx *TXOutput) String() string {
	return fmt.Sprintf("\n\t\t\tValue: %d, ScriptPubKey: %s", tx.Value, tx.ScriptPubKey)
}
