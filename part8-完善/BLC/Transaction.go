package BLC

import (
	"bytes"
	"log"
	"crypto/sha256"
	"fmt"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
	"time"
	"os"
)

type Transaction struct {
	TxID  []byte      //1. 交易hash
	Vins  []*TXInput  //2. 输入
	Vouts []*TXOutput //3. 输出
}

//1. 产生创世区块时的Transaction
func NewCoinbaseTransacion(address string) *Transaction {
	//创建创世区块交易的Vin
	txInput := &TXInput{[]byte{}, -1, nil, nil}
	//创建创世区块交易的Vout
	//txOutput := &TXOutput{10, address}
	txOutput := NewTxOutput(10, address)
	//生产交易Transaction
	txCoinBaseTransaction := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	//设置Transaction的TxHash
	txCoinBaseTransaction.SetID()

	return txCoinBaseTransaction

}

/*
	产生挖矿奖励交易Transaction  奖励为1
 */
func NewRewardTransacionYS() *Transaction {
	//创建创世区块交易的Vin
	txInput := &TXInput{[]byte{}, -1, nil, nil}
	//创建创世区块交易的Vout
	address := CoinbaseAddress(os.Getenv("NODE_ID"))
	if address == "" {
		//从钱包中取地址
		wallets := NewWallets(os.Getenv("NODE_ID"))
		for walletAddress, _ := range wallets.WalletMap {
			address = walletAddress
		}
	}

	if address== "" {
		log.Panic("未定义coinbase地址, 无法执行后续操作")
	}


	txOutput := NewTxOutput(1, address)
	//生产交易Transaction
	txCoinBaseTransaction := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	//设置Transaction的TxHash
	txCoinBaseTransaction.SetID()

	return txCoinBaseTransaction

}

//2. 创建普通交易产生的Transaction
func NewSimpleTransation(from string, to string, amount int64, utxoSet *UTXOSet, txs []*Transaction, nodeID string) *Transaction {
	//1.定义Input和Output的数组
	var txInputs []*TXInput
	var txOutputs []*TXOutput

	//获取本次转账要使用output
	//total, spentableUTXO := bc.FindSpentableUTXOs(from, amount, txs)
	total, spentableUTXO := utxoSet.FindSpentableUTXOs(from, amount, txs)

	//获取钱包的集合：
	wallets := NewWallets(nodeID)
	wallet := wallets.WalletMap[from]

	//判断本地钱包是否包含发送方公私钥
	if wallet == nil {
		log.Panic("本地钱包没有发送地址存档")
	}

	//2.创建Input
	for txID, indexArray := range spentableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TXInput{txIDBytes, index, nil, wallet.PublickKey}
			txInputs = append(txInputs, txInput)
		}
	}

	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput2 := NewTxOutput(total-amount, from)
	txOutputs = append(txOutputs, txOutput2)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	tx.SetID()
	//fmt.Println(tx)
	//设置签名
	utxoSet.blockChain.SignTrasanction(tx, wallet.PrivateKey, txs)

	return tx
}

func (tx *Transaction) IsCoinBaseTransaction() bool {
	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}

//签名
/*
签名：为了对一笔交易进行签名
	私钥：
	要获取交易的Input，引用的output，所在的之前的交易：
 */
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxsmap map[string]*Transaction) {
	//1.判断当前tx是否时coinbase交易
	if tx.IsCoinBaseTransaction() {
		return
	}

	//2.获取input对应的output所在的tx，如果不存在，无法进行签名
	for _, input := range tx.Vins {
		if prevTxsmap[hex.EncodeToString(input.TxID)] == nil {
			log.Panic("当前的Input，没有找到对应的output所在的Transaction，无法签名。。")
		}
	}

	//即将进行签名:私钥，要签名的数据
	txCopy := tx.TrimmedCopy()

	for index, input := range txCopy.Vins {
		// input--->5566

		prevTx := prevTxsmap[hex.EncodeToString(input.TxID)]

		txCopy.Vins[index].Signature = nil                                 //仅仅是一个双重保险，保证签名一定为空
		txCopy.Vins[index].PublicKey = prevTx.Vouts[input.Vout].PubKeyHash //设置input中的publickey为对应的output的公钥哈希

		txCopy.TxID = txCopy.NewTxID() //产生要签名的数据：

		//为了方便下一个input，将数据再置为空
		txCopy.Vins[index].PublicKey = nil

		//获取要交易的数据

		/*
		第一个参数
		第二个参数：私钥
		第三个参数：要签名的数据


		func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error)
		r + s--->sign
		input.Signatrue = sign
	 */
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxID)
		if err != nil {
			log.Panic(err)
		}

		sign := append(r.Bytes(), s.Bytes()...)
		tx.Vins[index].Signature = sign
	}

}

//获取要签名tx的副本
/*
要签名tx中，并不是所有的数据都要作为签名数据，生成签名
txCopy = tx{签名所需要的部分数据}
TxID

Inputs
	txid,vout,sign,publickey

Outputs
	value,pubkeyhash


交易的副本中包含的数据
	包含了原来tx中的输入和输出。
		输入中：sign，publickey
 */

func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs [] *TXInput
	var outputs [] *TXOutput
	for _, in := range tx.Vins {
		inputs = append(inputs, &TXInput{in.TxID, in.Vout, nil, nil})
	}

	for _, out := range tx.Vouts {
		outputs = append(outputs, &TXOutput{out.Value, out.PubKeyHash})
	}

	txCopy := &Transaction{[]byte{}, inputs, outputs}
	return txCopy

}

//将Transaction 序列化再进行 hash
func (tx *Transaction) SetID() {

	txBytes := gobEncode(tx)

	allBytes := bytes.Join([][]byte{txBytes, IntToHex(time.Now().Unix())}, []byte{})

	hash := sha256.Sum256(allBytes)
	//fmt.Printf("transationHash: %x", hash)
	tx.TxID = hash[:]
}

func (tx *Transaction) NewTxID() []byte {
	txCopy := tx
	//txCopy.TxID = []byte{}
	hash := sha256.Sum256(gobEncode(txCopy))
	return hash[:]
}

//验证交易
/*
验证的原理：
公钥 + 要签名的数据 验证 签名：rs
 */
func (tx *Transaction) Verifity(prevTxs map[string]*Transaction) bool {
	//判断当前input是否有对应的Transaction
	for _, input := range tx.Vins { //
		if prevTxs[hex.EncodeToString(input.TxID)] == nil {
			log.Panic("当前的input没有找到对应的Transaction，无法验证")
		}
	}

	//验证
	txCopy := tx.TrimmedCopy()

	curev := elliptic.P256() //曲线

	for index, input := range tx.Vins {
		//原理：再次获取 要签名的数据  + 公钥哈希 + 签名
		/*
		验证签名的有效性：
		第一个参数：公钥
		第二个参数：签名的数据
		第三、四个参数：签名：r，s
		func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool
		 */
		//ecdsa.Verify()

		//获取要签名的数据
		prevTx := prevTxs[hex.EncodeToString(input.TxID)]

		txCopy.Vins[index].Signature = nil
		txCopy.Vins[index].PublicKey = prevTx.Vouts[input.Vout].PubKeyHash
		txCopy.TxID = txCopy.NewTxID() //要签名的数据

		txCopy.Vins[index].PublicKey = nil

		//获取公钥
		/*
		type PublicKey struct {
			elliptic.Curve
			X, Y *big.Int
		}
		 */

		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.PublicKey)
		x.SetBytes(input.PublicKey[:keyLen/2])
		y.SetBytes(input.PublicKey[keyLen/2:])

		rawPublicKey := ecdsa.PublicKey{curev, &x, &y}

		//获取签名：

		r := big.Int{}
		s := big.Int{}

		signLen := len(input.Signature)
		r.SetBytes(input.Signature[:signLen/2])
		s.SetBytes(input.Signature[signLen/2:])

		if ecdsa.Verify(&rawPublicKey, txCopy.TxID, &r, &s) == false {
			fmt.Println("验证失败Verify")
			return false
		}

	}
	return true
}

//格式化输出
func (tx *Transaction) String() string {
	var vinStrings [][]byte
	for _, vin := range tx.Vins {
		vinString := fmt.Sprint(vin)
		vinStrings = append(vinStrings, []byte(vinString))
	}
	vinString := bytes.Join(vinStrings, []byte{})

	var outStrings [][]byte
	for _, out := range tx.Vouts {
		outString := fmt.Sprint(out)
		outStrings = append(outStrings, []byte(outString))
	}

	outString := bytes.Join(outStrings, []byte{})

	return fmt.Sprintf("\n\r\t\t===============================\n\r\t\tTxID: %x, \n\t\tVins: %v, \n\t\tVout: %v\n\t\t", tx.TxID, string(vinString), string(outString))
}
