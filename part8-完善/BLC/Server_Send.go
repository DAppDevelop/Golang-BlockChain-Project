package BLC

import (
	"net"
	"log"
	"io"
	"bytes"
)

/*
	所有消息都是通过这个方法来发送到其他节点
 */
func sendData(to string, data []byte) {
	//fmt.Println("向",to,"发送",data)
	conn, err := net.Dial("tcp", to)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	//发送数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

/*
	发送本地版本/区块高度
 */
func sendVersion(to string, bc *Blockchain) {
	//1.创建对象
	bestHeight := bc.GetBestHeight()//获取当前节点区块链高度
	version := &Version{NODE_VERSION, bestHeight, nodeAddress}

	sendCommandData(COMMAND_VERSION, version, to)
}

/*
	发送请求要获取对方blockhash的消息
 */
func sendGetBlocksHash(to string) {
	//1.创建对象
	getBlocks := GetBlocks{nodeAddress}

	sendCommandData(COMMAND_GETBLOCKS, getBlocks, to)
}

/*
	发送所有blockHash 数组的消息
 */
func sendInv(to string, kind string, data [][]byte) {
	//1.创建对象
	inv := Inv{nodeAddress, kind, data}

	sendCommandData(COMMAND_INV, inv, to)
}

/*
	发送请求对方根据hash返回对应的block的消息
 */
func sendGetData(to string, kind string, hash []byte) {
	//1.创建对象
	getData := GetData{nodeAddress, kind, hash}

	sendCommandData(COMMAND_GETDATA, getData, to)
}

/*
	发送block对象给对方
 */
func sendBlock(to string, block *Block) {
	//1.创建对象
	blockData := BlockData{nodeAddress, gobEncode(block)}

	sendCommandData(COMMAND_BLOCKDATA, blockData, to)
}

/*
	发送交易信息到主节点
 */
func sendTransactionToMainNode(to string, txs []*Transaction)  {
	sendCommandData(COMMAND_TXS, txs, to)
}

func sendTransactionToMiner(to string, txs []*Transaction)  {
	sendCommandData(COMMAND_REQUIREMINE, txs, to)
}



func sendCommandData(command string, data interface{}, to string)  {
	//2.对象序列化为[]byte
	payload := gobEncode(data)
	//3.拼接命令和对象序列化
	request := append(commandToBytes(command), payload...)
	//4.发送消息
	sendData(to, request)
}


