package BLC

import "time"

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
func NewBlock(data []byte, height int64, preHash []byte) *Block {
	return &Block{
		BlockHeight: height,
		PreHash:     preHash,
		Data:        data,
		Timestamp:   time.Now().Unix(),
		Hash:        nil,	//TODO
	}
}