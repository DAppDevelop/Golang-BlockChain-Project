package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

type ProofOfWork struct {
	Block *Block // 当前要验证的区块

	//对于整数的高精度计算 Go 语言中提供了 big 包。有用来表示大整数的 big.Int
	target *big.Int // 当hash小于此target时，为挖矿成功
}

func NewProofOfWork(block *Block) *ProofOfWork {
	//1. 创建一个初始值为1的target
	target := big.NewInt(1)

	//2. 左移 256 - targetBit
	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) Run() ([]byte, int64) {
	//使用nonce计算hash不符合target时候，加1，直到hash符合要求
	nonce := 0

	var hashInt big.Int
	var hash [32]byte
	dataBytes := pow.prepareData()
	for {
		//1. 将Block的属性拼接成字节数组作为sha256.Sum256的入参
		dataBytes := bytes.Join(
			[][]byte{ //[]byte的切片
				dataBytes,
				IntToHex(int64(nonce)),
			},
			[]byte{},
		)

		//2. 生成hash
		hash = sha256.Sum256(dataBytes)
		//fmt.Printf("\r%x", hash)

		//将hash转换成*int类型并返回给hashInt
		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于Block里面的target

		//3. 判断hash有效性，如果满足条件，跳出循环
		if pow.target.Cmp(&hashInt) == 1 {
			fmt.Printf("\nhash: %x\n", hash) //hash: 00ea9e3743900b6086acbb86390457f72fb3a4908609bd900536064f8e89448d
			break
		}

		//如果不满足条件，nonce+1并继续循环
		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}



// 数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData() []byte {
	//bytes.Join 以sep为连接符，拼接[][]byte
	//提取这个地方的数据, 只改变nonce,其他不用重新运算
	data := bytes.Join([][]byte{ //[]byte的切片
		pow.Block.PrevBlockHash,
		pow.Block.HashTransactions(),
		IntToHex(pow.Block.Timestamp),
		IntToHex(int64(targetBit)),
		IntToHex(int64(pow.Block.Height)),
	}, []byte{},
	)

	return data
}

func (pow *ProofOfWork) IsValid() bool {

	var hashInt big.Int

	hashInt.SetBytes(pow.Block.Hash)

	//1.proofOfWork.Block.Hash
	//2.proofOfWork.Target 作比较
	if pow.target.Cmp(&hashInt) == 1 {
		return true
	}

	return false
}
