package BLC

import (
	"net"
	"log"
	"io"
	"bytes"
	"fmt"
)

/*
	所有消息都是通过这个方法来发送到其他节点
 */
func sendData(to string, data []byte)  {
	fmt.Println("向",to,"发送",data)
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

func sendVersion(to string, bc *Blockchain)  {
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