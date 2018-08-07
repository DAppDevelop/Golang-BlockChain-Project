package BLC

type TXOutput struct {
	Value        int64  //金额
	ScriptPubKey string //用户名(scriptPubkey:锁定脚本,包含公钥)
}
