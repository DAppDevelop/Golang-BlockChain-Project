package BLC

import (
	"fmt"
	"encoding/hex"
)

type UTXO struct {
	TxID   []byte    //1.该Output所在的交易id
	Index  int       //2.该Output 的下标
	Output *TXOutput //3.Output
}


//打印格式
func (utxo *UTXO) String() string {
	return fmt.Sprintf(
		"\n------------------------------"+
			"\nA UTXO's Info:\n\t"+
			"TxID:%s,\n\t"+
			"Index:%d,\n\t"+
			"Output: %v,\n\t",
		hex.EncodeToString(utxo.TxID),
		utxo.Index,
		utxo.Output,
		)
}
