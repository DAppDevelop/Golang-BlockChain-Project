package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	//1. 区块高度
	Height int64
	//2. 上一个区块HASH
	PrevBlockHash []byte
	//3. 交易数据
	Data []byte
	//4. 时间戳
	Timestamp int64
	//5. Hash
	Hash []byte
	//6. Nonce
	Nonce int64
}

func NewBlock(data string, height int64, preBlockHash []byte) *Block {
	block := &Block{
		height,
		preBlockHash,
		[]byte(data),
		time.Now().Unix(),
		nil,
		0}

	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := NewProofOfWork(block)

	// 挖矿验证
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	//block.SetHash()

	return block

}

//func (block *Block) SetHash() {
//	//1.转换Height  256进制
//	heightBytes := IntToHex(block.Height)
//
//	//fmt.Printf("heightBytes: %v %T\n", heightBytes, heightBytes)
//
//	//2.转换时间戳  转成二进制再转[]byte
//	timeString := strconv.FormatInt(block.Timestamp, 2)
//	//fmt.Println(block.Timestamp)
//	//fmt.Printf("timeString: %v %T\n", timeString, timeString)
//
//	timeBytes := []byte(timeString)
//	//fmt.Printf("timeBytes: %v %T\n", timeBytes, timeBytes)
//	//  3. 拼接所有属性
//	blockBytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timeBytes}, []byte{})
//
//	//fmt.Printf("blockBytes: %v %T\n", blockBytes, blockBytes)
//	// 4. 生成Hash
//	hash := sha256.Sum256(blockBytes)
//	//fmt.Printf("hash: %v %T\n", hash, hash)
//
//	block.Hash = hash[:]
//}

func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func (block *Block) String() string {
	return fmt.Sprintf(
		"\n------------------------------"+
			"\nABlock's Info:\n\t"+
			"Height:%d,\n\t"+
			"PreHash:%v,\n\t"+
			"Data: %v,\n\t"+
			"Timestamp: %s,\n\t"+
			"Hash: %x,\n\t"+
			"Nonce: %v\n\t",
		block.Height,
		block.PrevBlockHash,
		block.Data,
		time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"),
		block.Hash, block.Nonce)
}

// 将区块序列化成字节数组
func (block *Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func DeserializeBlock(blockBytes []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}