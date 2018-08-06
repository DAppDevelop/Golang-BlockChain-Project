package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

//0000 0000 0000 0000 1001 0001 0000 .... 0001
// 256位Hash里面前面至少要有16个零
const targetBit = 16

type ProofOfWork struct {
	Block  *Block   // 当前要验证的区块
	target *big.Int // 大数据存储 2^24
}

func NewProofOfWork(block *Block) *ProofOfWork {
	//1. 创建一个初始值为1的target
	target := big.NewInt(1)

	//2. 左移256 - targetBit
	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{block, target}
}

func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
	nonce := 0

	var hashInt big.Int
	var hash [32]byte

	for {
		//1. 将Block的属性拼接成字节数组
		dataBytes := proofOfWork.prepareData(nonce)

		//2. 生成hash
		hash = sha256.Sum256(dataBytes)
		//fmt.Printf("\r%x", hash)


		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于Block里面的target
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//3. 判断hash有效性，如果满足条件，跳出循环
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			fmt.Printf("hash: %x\n", hash)
			break
		}

		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}



// 数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *ProofOfWork) IsValid() bool {

	//1.proofOfWork.Block.Hash
	//2.proofOfWork.Target
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)


	// Cmp compares x and y and returns:
	//
	//   -1 if x <  y
	//    0 if x == y
	//   +1 if x >  y
	if proofOfWork.target.Cmp(&hashInt) == 1 {
		return true
	}

	return false
}
