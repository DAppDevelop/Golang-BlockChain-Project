package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey //1.私钥
	PublickKey []byte           //2.公钥  原始公钥
}

/*
	创建新的钱包
 */
func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()

	return &Wallet{privateKey, publicKey}
}

//产生一对密钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	/*
	1.根据椭圆曲线算法，产生随机私钥
	2.根据私钥，产生公钥
	椭圆：ellipse，
	曲线：curve，

	椭圆曲线加密：(ECC：ellipse curve Cryptography)，非对称加密
		加密：
			对称加密和非对称机密啊

		SECP256K1,算法

		x轴(32byte)，y轴(32byte)--->

	 */
	//椭圆加密

	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey
}

//根据公钥获取对应的地址
func (w *Wallet) GetAddress() []byte {
	/*
	1.原始公钥-->sha256-->160-->公钥哈希
	2.版本号+公钥哈希--->校验码
	3.版本号+公钥哈希+校验码--->Base58编码
	 */

	//step1：得到公钥哈希
	pubKeyHash := PubKeyHash(w.PublickKey)

	//step2：添加版本号：
	versioned_payload := append([]byte{version}, pubKeyHash...)

	//step3：根据versioned_payload-->两次sha256,取前4位，得到checkSum
	checkSumBytes := CheckSum(versioned_payload)
	//step4：拼接全部数据
	full_payload := append(versioned_payload, checkSumBytes...)
	//step5：Base58编码
	address := Base58Encode(full_payload)

	return address
}

/*
原始公钥-->公钥哈希
1.sha256
2.ripemd160
 */
func PubKeyHash(publickKey []byte) []byte {
	//1.sha256
	hasher := sha256.New()
	hasher.Write(publickKey)
	hash1 := hasher.Sum(nil)
	//2.ripemd160

	hasher2 := ripemd160.New()
	hasher2.Write(hash1)
	hash2 := hasher.Sum(nil)
	//3.返回

	return hash2
}

//产生校验码
/*
两次sha256
 */
func CheckSum(payload [] byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	//取前四个字节
	return secondHash[:addressCheckSumLen]

}

//校验地址是否有效：
func IsValidAddress(address []byte) bool {
	//step1：Base58解码
	//version+pubkeyHash+checksum
	full_payload := Base58Decode(address)

	//step2：获取地址中携带的checkSUm
	checkSumBytes := full_payload[len(full_payload)-addressCheckSumLen:]
	version_payload := full_payload[:len(full_payload)-addressCheckSumLen]
	//step3：versioned_payload，生成一次校验码
	checkSumBytes2 := CheckSum(version_payload)
	//step4：比较checkSumBytes，checkSumBytes2
	return bytes.Compare(checkSumBytes, checkSumBytes2) == 0
}
