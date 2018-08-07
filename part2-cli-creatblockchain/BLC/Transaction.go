package BLC

type Transaction struct {
	TxHash []byte      //1. 交易hash
	Vins   []*TXInput  //2. 输入
	Vouts  []*TXOutput //3. 输出
}
