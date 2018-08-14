package BLC

type UTXO struct {
	TxID   []byte    //1.该Output所在的交易id
	Index  int       //2.该Output 的下标
	Output *TXOutput //3.Output
}
