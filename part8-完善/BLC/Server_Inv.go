package BLC

type Inv struct {
	AddrFrom string   //当前节点地址
	Type     string   //类型（block or Transaction
	Items    [][]byte //对应类型的数据的hash
}
