package blc

import (
	"github.com/boltdb/bolt"
	"log"
	"simple_bitcoin/utils"
)

// 区块链迭代器
type BlockChainIterator struct {
	 currentHash []byte
	 db *bolt.DB
}

// 获取当前hash区块，递进到下一个区块
func (bci *BlockChainIterator) Next() (block *Block) {
	var err error
	err = bci.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BucketName))
		// 到达创世块停止
		if string(bci.currentHash) == string(INITIALLY_HASH) {
			return nil
		}
		blockBytes := bucket.Get(bci.currentHash)
		block = Deserialization(blockBytes)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	// 更新迭代器hash
	if string(bci.currentHash) != string(INITIALLY_HASH) {
		bci.currentHash = block.PreHash
	}
	return
}