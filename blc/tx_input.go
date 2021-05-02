package blc

import (
	"bytes"
	"fmt"
	"strings"
)

type TxInput struct {
	// 与上一个交易相关的参数，指明交易的来源
	OutputTxHash []byte
	OutputIdx int
	// 签名与公钥
	Signature []byte
	PubKey []byte
}

// 比对公钥的Hash即PubKeyHash
// 输入存储的是原生的公钥所以要先hash处理（SHA256 + RIPEMD160）
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)			// 对公钥取Hash
	return bytes.Compare(pubKeyHash, lockingHash) == 0
}

func (in *TxInput) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("OutputTxHash : %x\n", in.OutputTxHash))
	lines = append(lines, fmt.Sprintf("OutputIdx : %d\n", in.OutputIdx))
	lines = append(lines, fmt.Sprintf("Signature : %x\n", in.Signature))
	lines = append(lines, fmt.Sprintf("PubKey : %x\n", in.PubKey))
	return strings.Join(lines, "")
}