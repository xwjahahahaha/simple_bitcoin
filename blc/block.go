package blc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"time"
)

//创世块初始Hash
var INITIALLY_HASH = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type Block struct {
	//1.区块高度
	BlockHeight int64
	//2.上一个区块hash
	PreHash []byte
	//3.交易数据
	Transactions []*Transaction
	//4.时间戳
	Timestamp int64
	//5.本区块的Hash
	Hash []byte
	//6.nonce随机值
	Nonce int64
}

//创建新的区块
func NewBlock(transaction []*Transaction, height int64, preHash []byte) *Block {
	newBlock := &Block{
		BlockHeight: height,
		PreHash:     preHash,
		Transactions: transaction,
		Timestamp:   time.Now().Unix(),
		Hash:        nil,
		Nonce:  	 0,
	}
	//运行工作量证明计算nonce
	//为当前区块生成工作量证明对象
	pow := NewProofOfWork(newBlock)
	//计算
	nonce, hash := pow.Run()
	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce

	return newBlock
}

// 生成创世区块
func CreateGenesisBlcok(coinBase *Transaction, height int64) *Block {
	return NewBlock([]*Transaction{coinBase}, height, INITIALLY_HASH)
}


// 区块序列化
func (block *Block) Serialization() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(block)
	if err != nil {
		fmt.Println("Serialization err : ", err)
	}
	return res.Bytes()
}

// 区块反序列化
func Deserialization(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Deserialization err : ", err)
	}
	return &block
}

// 区块信息打印
func (block *Block) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("================================================================================================\n"))
	lines = append(lines, fmt.Sprintf("区块高度: %d\n", block.BlockHeight))
	lines = append(lines, fmt.Sprintf("前区块hash值：%x\n", block.PreHash))
	lines = append(lines, fmt.Sprintf("本区块Hash值: %x\n", block.Hash))
	for _, tx := range block.Transactions {
		lines = append(lines, fmt.Sprintln("交易数据：\n", tx))
	}
	//转换下时间
	timeFormat := time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05")
	lines = append(lines, fmt.Sprintf("本区块时间：%s\n", timeFormat))
	lines = append(lines, fmt.Sprintf("本区块随机数Nonce：%d\n", block.Nonce))
	lines = append(lines, fmt.Sprintf("================================================================================================\n"))
	return strings.Join(lines, "\n")
}

// 序列化区块所有的交易（方便计算区块Hash以及Nonce）
func (block *Block) SerializeAllTxs() []byte {
	var res []byte
	for _, tx := range block.Transactions {
		res = append(res , tx.Serialize()...)
	}
	return res
}