package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
	return NewBlock(data, 0, INITIALLY_HASH)
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
func (block *Block)BlockChainStdOutPrint()  {
	// 输出，这里只是标准stdout，可以改为其他输出
	fmt.Printf("================================================================================================\n")
	fmt.Printf("前区块hash值：%x\n", block.BlockHeight)
	fmt.Printf("本区块hash值：%x\n", block.PreHash)
	fmt.Printf("梅克尔根：%x\n", block.Data)
	//转换下时间
	timeFormat := time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05")
	fmt.Printf("本区块时间：%s\n", timeFormat)
	fmt.Printf("本区块随机数Nonce：%d\n", block.Nonce)
	fmt.Printf("================================================================================================\n")
}