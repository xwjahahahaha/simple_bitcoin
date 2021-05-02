package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
	"simple_bitcoin/utils"
)

type Wallet struct {
	// 公私钥对
	PrivateKey ecdsa.PrivateKey
	Publickey []byte
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
	versionByte := utils.Int64ToBytes(int64(utils.PubKeyVersion))
	versionedPayload := append(versionByte[len(versionByte)-1:], publicSHA256...)
	checkSum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checkSum...)
	address := Base58Encode(fullPayload)
	return address
}

// 对公钥先SHA256再RIPEMD160
func HashPubKey(pubKey []byte) []byte {
	// 先sha256再RIPEMD160加密
	pubKeySha256 := sha256.Sum256(pubKey)
	// 创建RIPEMD160
	RIPEMD160Hasher := ripemd160.New()			// 注意这里用golang.org/x/crypto/ripemd160包下的ripemd160
	_, err := RIPEMD160Hasher.Write(pubKeySha256[:])
	if err != nil {
		log.Panic(err)
	}
	return RIPEMD160Hasher.Sum(nil)
}

/**
 * @Description:  计算checkSum
 * @param payload
 * @return []byte
 */
func checksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:utils.AddressCheckSumLen]	// 截取hash后AddressCheckSumLen个字节
}


/**
 * @Description: 通过checksum验证地址是否有效（符合生成规则checksum而不是随机乱码）
 * @param address
 * @return bool
 */
func ValidateAddress(address string) bool  {
	// 1. base58解码
	pubKeyHash := Base58Decode([]byte(address))
	// 2. 获取checkSum和pubHashKey部分
	// prefixLoad [1:] 是去掉Base58 Decode前面添加的一个字节0
	prefixLoad, actualChecksum := pubKeyHash[1:len(pubKeyHash)-utils.AddressCheckSumLen], pubKeyHash[len(pubKeyHash)-utils.AddressCheckSumLen:]
	// 3. 计算sum
	nowChecksum := checksum(prefixLoad)
	//fmt.Printf("prefixLoad : %x\n ", prefixLoad)
	return bytes.Compare(nowChecksum, actualChecksum) == 0
}

/**
 * @Description:  解析地址为PubKeyHash
 * @param address
 * @return []byte
 */
func ResolveAddressToPubKeyHash(address string) []byte {
	if ! ValidateAddress(address) {
		log.Panic("InValid Address")
		return nil
	}
	// base58解码
	pubKeyHash := Base58Decode([]byte(address))
	// 2byte 是因为 1byte由Decode造成前面1byte的0，剩下1byte是version
	return pubKeyHash[2:len(pubKeyHash)-utils.AddressCheckSumLen]
}
