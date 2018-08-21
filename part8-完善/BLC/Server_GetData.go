package BLC

type GetData struct {
	AddrFrom string //当前节点自己的地址
	Type string //数据类型（block或者tx)
	Hash []byte//block或者Tx的hash
}
