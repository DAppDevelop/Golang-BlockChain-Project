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
	//获取当前节点区块链高度
	bestHeight := bc.GetBestHeight()

	//创建Version对象(为version命令消息要传送的数据)
	version := &Version{NODE_VERSION, bestHeight, nodeAddress}

	//序列化version对象
	payload := gobEncode(version)

	//将序列化的version命令和数据payload拼接成特定格式的[]byte
	request := append(commandToBytes(COMMAND_VERSION), payload...)

	//发送
	sendData(to, request)

}

/*
	发送请求要获取对方blockhash的消息
 */
func sendGetBlocksHash(to string) {
	//1.创建对象
	getBlocks := GetBlocks{nodeAddress}
	//2.对象序列化为[]byte
	payload := gobEncode(getBlocks)
	//3.拼接命令和对象序列化
	request := append(commandToBytes(COMMAND_GETBLOCKS), payload...)
	//4.发送消息
	sendData(to, request)
}

/*
	发送所有blockHash 数组的消息
 */
func sendInv(to string, kind string, data [][]byte) {
	//1.创建对象
	inv := Inv{nodeAddress, kind, data}
	//2.对象序列化为[]byte
	payload := gobEncode(inv)
	//3.拼接命令和对象序列化
	request := append(commandToBytes(COMMAND_INV), payload...)
	//4.发送消息
	sendData(to, request)
}

/*

 */
func sendGetData(to string, kind string, hash []byte) {
	//1.创建对象
	getData := GetData{nodeAddress, kind, hash}
	//2.对象序列化为[]byte
	payload := gobEncode(getData)
	//3.拼接命令和对象序列化
	request := append(commandToBytes(COMMAND_GETDATA), payload...)
	//4.发送消息
	sendData(to, request)
}

func sendBlock(to string, block *Block) {
	//1.创建对象
	blockData := BlockData{nodeAddress, block.Serialize()}
	//2.对象序列化为[]byte
	payload := gobEncode(blockData)
	//3.拼接命令和对象序列化
	request := append(commandToBytes(COMMAND_BLOCKDATA), payload...)
	//4.发送消息
	sendData(to, request)
}
