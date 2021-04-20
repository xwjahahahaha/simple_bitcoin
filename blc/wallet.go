package blc

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"simple_bitcoin/utils"
)

type Wallet struct {
	// 公私钥对
	PrivateKey ecdsa.PrivateKey
	Publickey []byte
}

// 多个钱包
type Wallets struct {
	Wallets Wallet
}

// 新建一个钱包
func NewWallet() *Wallet {
	private, public := newKeyPair()		// 创建公私钥
	return &Wallet{
		PrivateKey: private,
		Publickey:  public,
	}
}

// 创建公私钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()		// 创建曲线
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)	// 创建私钥
	if err != nil {
		log.Panic(err)
	}
	// 创建公钥
	// 基于椭圆曲线，公钥是曲线上的点，所以公钥是X，Y坐标的组合, 在比特币中将其连接起来实现公钥
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

// 获取钱包地址
func (w Wallet) GetAddress() []byte {
	publicSHA256 := HashPubKey(w.Publickey)
	versionedPayload := append([]byte(utils.PubKeyVersion), publicSHA256...)
	checkSum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checkSum...)
	address := Base58Encode(fullPayload)
	return address
}

func HashPubKey(pubKey []byte) []byte {
	// 先sha256再RIPEMD160加密
	pubKeySha256 := sha256.Sum256(pubKey)
	// 创建RIPEMD160
	RIPEMD160Hasher := crypto.RIPEMD160.New()
	_, err := RIPEMD160Hasher.Write(pubKeySha256[:])
	if err != nil {
		log.Panic(err)
	}
	return RIPEMD160Hasher.Sum(nil)
}

func checksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:utils.AddressCheckSumLen]	// 截取hash后AddressCheckSumLen个字节
}