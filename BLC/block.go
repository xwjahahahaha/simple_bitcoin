package BLC

import (
	"time"
)

//创世块初始Hash
var INIIALLY_HASH = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type Block struct {
	//1.区块高度
	BlockHeight int64
	//2.上一个区块hash
	PreHash []byte
	//3.交易数据
	Data []byte
	//4.时间戳
	Timestamp int64
	//5.本区块的Hash
	Hash []byte
	//6.nonce随机值
	Nonce int64
}

//创建新的区块
func NewBlock(data string, height int64, preHash []byte) *Block {
	newBlock := &Block{
		BlockHeight: height,
		PreHash:     preHash,
		Data:        []byte(data),
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
func CreateGenesisBlcok(data string) *Block {
	return NewBlock(data, 0, INIIALLY_HASH)
}
