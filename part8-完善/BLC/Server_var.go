package BLC

var knowNodes = []string{"localhost:3000","localhost:3002","localhost:3001"}//主节点地址/挖矿节点/普通节点

var nodeAddress string//当前节点地址

var blockArray [][]byte //记录尚未同步的区块的hash

var coinbaseAddress string //挖矿奖励分配地址