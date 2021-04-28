package blc

import "bytes"

type TxInput struct {
	// 与上一个交易相关的参数
	OutputTxHash []byte
	OutputIdx int
	// 签名与公钥
	Signature []byte
	PubKey []byte
}

// 比对公钥的Hash即PubKeyHash
// 输入存储的是原生的公钥所以要先hash处理（SHA256 + RIPEMD160）
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)	// 对公钥取Hash
	return bytes.Compare(pubKeyHash, lockingHash) == 0
}

