package BLC

import (
	"encoding/gob"
	"bytes"
	"log"
)

func handleVersion(request []byte, bc *Blockchain) {

	//1.从request中获取版本的数据：[]byte
	versionBytes := request[COMMAND_LENGTH:]

	//2.反序列化--->version
	version := &Version{}

	decoder := gob.NewDecoder(bytes.NewReader(versionBytes))

	err := decoder.Decode(&version)
	if err != nil {
		log.Panic(err)
	}

	//3.操作bc，获取自己的最后block的height
	height := bc.GetBestHeight()
	foreignerBestHeight := version.BestHeight

	//4.根对方的比较, 相同则不做操作
	if height > foreignerBestHeight {
		//当前节点比对方节点高度高
		sendVersion(version.AddrFrom, bc)
	} else if foreignerBestHeight > height {
		//当前节点比对方节点高度低,向对方节点请求对方节点的blockchain hash集
		sendGetBlocksHash(version.AddrFrom)
	}

}

func handleGetBlocksHash(request []byte, bc *Blockchain)  {

}

func handleInv(request []byte, bc *Blockchain)  {

}

func handleGetData(request []byte, bc *Blockchain)  {

}

func handleGetBlockData(request []byte, bc *Blockchain)  {

}


