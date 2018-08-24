package BLC

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
)

/*
	启动服务器
 */
func startServer(nodeID string, mineAddress string) {
	//设置coinbase
	coinbaseAddress = mineAddress
	//拼接nodeID到ip后
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	//监听地址
	listener, err := net.Listen("tcp", nodeAddress)

	if err != nil {
		log.Panic(err)
	}

	defer listener.Close()

	bc := BlockchainObject(nodeID)
	//defer bc.DB.Close()

	//判断是否为主节点, 非主节点的节点需要向主节点发送Version消息
	//fmt.Println(nodeAddress, knowNodes[0])
	if nodeAddress != knowNodes[0] {
		//fmt.Println("sendVersion")
		sendVersion(knowNodes[0], bc)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Panic(err)
		}

		fmt.Println("发送方已接入..", conn.RemoteAddr())

		go handleConnection(conn, bc)
	}
}

/*
	处理请求结果
 */
func handleConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(request[:COMMAND_LENGTH])

	fmt.Printf("接收到的命令是：%s\n", command)

	switch command {
	case COMMAND_VERSION:
		handleVersion(request, bc)
	case COMMAND_GETBLOCKS:
		handleGetBlocksHash(request, bc)
	case COMMAND_INV:
		handleInv(request, bc)
	case COMMAND_GETDATA:
		handleGetData(request, bc)
	case COMMAND_BLOCKDATA:
		handleGetBlockData(request, bc)
	case COMMAND_TXS:
		handleTransactions(request, bc)
	case COMMAND_REQUIREMINE:
		handleRequireMine(request, bc)
	case COMMAND_VERIFYBLOCK:
		handleVerifyBlock(request, bc)
	default:
		fmt.Println("无法识别....")
	}

	defer conn.Close()
}
