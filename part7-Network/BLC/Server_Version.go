package BLC

//定义为12字节长度
type Version struct {
	Version    int64  //版本
	BestHeight int64  //当前节点区块链中最后一个区块的高度
	AddrFrom   string //当前节点地址
}
