package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Height        int64  //1. 区块高度
	PrevBlockHash []byte //2. 上一个区块HASH
	Data          []byte //3. 交易数据
	Timestamp     int64  //4. 时间戳
	Hash          []byte //5. Hash
	Nonce         int64  //6. Nonce
}

func NewBlock(data string, height int64, preBlockHash []byte) *Block {
	block := &Block{
		height,
		preBlockHash,
		[]byte(data),
		time.Now().Unix(),
		nil,
		0}

	//创建工作量证明结构体
	pow := NewProofOfWork(block)

	//调用工作量证明的方法并且返回有效的Hash和Nonce（挖矿）
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

//打印格式
func (block *Block) String() string {
	return fmt.Sprintf(
		"\n------------------------------"+
			"\nABlock's Info:\n\t"+
			"Height:%d,\n\t"+
			"PreHash:%x,\n\t"+
			"Data: %s,\n\t"+
			"Timestamp: %s,\n\t"+
			"Hash: %x,\n\t"+
			"Nonce: %v\n\t",
		block.Height,
		block.PrevBlockHash,
		block.Data,
		time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"),
		block.Hash, block.Nonce)
}

// 序列化：将区块序列化成字节数组
func (block *Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(result.Bytes())
	return result.Bytes()
}

// 反序列化：将字节数组反序列化为block对象
func DeserializeBlock(blockBytes []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
