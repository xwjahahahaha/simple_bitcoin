package blc

import (
	"bytes"
	"fmt"
	"simple_bitcoin/utils"
	"strings"
)

type TxOutput struct {
	// 值
	Value int
	// 解锁规则
	PubKeyHash []byte		// 公钥的Hash
}

// 给当前地址的输出上锁（设置PubKeyHash）
func (out *TxOutput) Lock(address []byte){
	// 1. 将地址解码base58
	pubKeyHash := Base58Decode(address)
	// 2. 截取中间段就是PubKeyHash
	// 前一个是0，后一个byte是version，后四个是checksum
	pubKeyHash = pubKeyHash[2:len(pubKeyHash)-utils.AddressCheckSumLen]
	fmt.Printf("address = %s, output pubKeyHash = %x\n", address, pubKeyHash)
	// 3. 设置
	out.PubKeyHash = pubKeyHash
}

// 检查此公钥Hash是否被用于锁定输出
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(pubKeyHash, out.PubKeyHash) == 0
}

/**
 * @Description: 创建一个vout
 * @param value
 * @param address
 * @return *TxOutput
 */
func NewTxOutput(value int, address string) *TxOutput {
	newTxout :=  &TxOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	newTxout.Lock([]byte(address)) 	// 上锁（地址 => PubkeyHash）
	return newTxout
}


func (out *TxOutput) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Value : %d\n", out.Value))
	lines = append(lines, fmt.Sprintf("PubKeyHash : %x\n", out.PubKeyHash))
	return strings.Join(lines, "")
}