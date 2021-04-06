package blc

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"simple_bitcoin/utils"
)

var blockHeightCount int64

//区块链数据结构
type BlockChain struct {
	//所有的区块
	LastHash []byte
	DB *bolt.DB
}

// 生成区块链并持久化存储
func CreateBlockchainDB(nodeID string) *BlockChain  {
	// 拼接数据库文件名字符串
	dbName := utils.DBName + "_" + nodeID + ".db"
	// 检查该文件是否已存在
	if dbExists(dbName){
		// 已存在就退出
		fmt.Println(dbName, " Blockchain already exists.")
		os.Exit(1)
	}
	// 创建创世区块
	GenesisBlock := CreateGenesisBlcok(utils.GenesisCoinbaseData)
	// 打开数据库，如果文件不存在会自动创建
	db, err := bolt.Open(dbName, 0666, nil)
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

// 创建新区块链实例（读取区块链数据库文件创建区块链实例）
// 注意：必须有区块链数据文件生成才可
func NewBlockchain(nodeID string) *BlockChain {
	// 拼接数据库文件名字符串
	dbName := utils.DBName + "_" + nodeID + ".db"
	// 检查该文件是否已存在
	if !dbExists(dbName){
		// 不存在需要创建
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	// 打开数据库
	db, err := bolt.Open(dbName, 0666, nil)
	if err != nil {
		log.Panic(err)
	}
	var lastHash []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket:= tx.Bucket([]byte(utils.BucketName))
		// 获取当前数据库最后的hash
		lastHash = bucket.Get([]byte(utils.LastHashKey))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//添加到区块链中
	return &BlockChain{
		LastHash: lastHash,
		DB:       db,
	}
}

// 区块链中添加新的区块
func (bc *BlockChain) AddNewBlock(data string)  {
	var err error
	blockHeightCount ++
	//计算新区块数据
	newHeight, preHash := blockHeightCount,	bc.LastHash
	//生成新区块
	newBlock := NewBlock(data, newHeight, preHash)
	//加入区块链
	err = bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BucketName))
		// 存储新区块
		err = bucket.Put(newBlock.Hash, newBlock.Serialization())
		if err != nil {
			log.Panic(err)
		}
		// 更新最后的hash
		err = bucket.Put([]byte(utils.LastHashKey), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		// 更新区块链记录的lastHash
		bc.LastHash = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 输出区块链(stdout打印)
func (bc *BlockChain) PrintBlockChain()  {
	iterator := &BlockChainIterator{
		currentHash: bc.LastHash,
		db:          bc.DB,
	}
	block := iterator.Next()
	for block != nil {
		block.BlockChainStdOutPrint()
		block = iterator.Next()
	}
}

// 判断数据库是否存在
func dbExists(dbFile string) bool {
	// os.stat判断文件是否存在
	if _, err := os.Stat(dbFile); os.IsNotExist(err){
		return false
	}
	return true
}