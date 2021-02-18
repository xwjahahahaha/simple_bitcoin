package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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
}

//1.创建新的区块
func NewBlock(data string, height int64, preHash []byte) *Block {
	newBlock := &Block{
		BlockHeight: height,
		PreHash:     preHash,
		Data:        []byte(data),
		Timestamp:   time.Now().Unix(),
		Hash:        nil,
	}
	newBlock.SetBlockHash()
	return newBlock
}

//2.计算当前区块的Hash
func (block *Block) SetBlockHash() {
	//数据类型统一为[]byte
	height := IntToBytes(block.BlockHeight)
	//时间戳 => 二进制字符串 => []byte
	timeStamp := []byte(strconv.FormatInt(block.Timestamp, 2))
	//合并为二位字符数组切片
	blockData := bytes.Join([][]byte{height, timeStamp, block.PreHash, block.Data}, []byte{})
	//求hash,返回256位/32位字节数组
	hashAry := sha256.Sum256(blockData)
	//转换为字节数组切片
	block.Hash = hashAry[:]
}

//3.生成创世区块
func CreateGenesisBlcok(data string) *Block {
	return NewBlock(data, 0, INIIALLY_HASH)
}
