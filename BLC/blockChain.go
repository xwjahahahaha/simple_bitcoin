package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"simple_bitcoin/utils"
)

//区块链数据结构
type BlockChain struct {
	//所有的区块
	LastHash []byte
	DB *bolt.DB
}

//创建新区块链
func CreateNewBC() *BlockChain {
	//创建创世区块
	GenesisBlock := CreateGenesisBlcok(utils.GenesisString)
	// 打开数据库
	db, err := bolt.Open(utils.DBName, 0666, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, bucketErr := tx.CreateBucket([]byte(utils.BucketName))
		if bucketErr != nil {
			log.Panic(bucketErr)
		}
		// 存储创世块
		putErr := bucket.Put(GenesisBlock.Hash, GenesisBlock.Serialization())
		if putErr != nil {
			log.Panic(putErr)
		}
		// 存储最后的Hash
		putErr = bucket.Put([]byte(utils.LastHashKey), GenesisBlock.Hash)
		if putErr != nil {
			log.Panic(putErr)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//添加到区块链中
	return &BlockChain{
		LastHash: GenesisBlock.Hash,
		DB:       db,
	}
}

//区块链中添加新的区块
//func (bc *BlockChain) AddNewBlock(data string)  {
//	length := len(bc.Blocks)
//	lastBlock := bc.Blocks[length-1]
//	//计算新区块数据
//	newHeight, preHash := int64(length), lastBlock.Hash
//	//生成新区块
//	newBlock := NewBlock(data, newHeight, preHash)
//	//加入区块链
//	bc.Blocks = append(bc.Blocks, newBlock)
//}

