package blc

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"simple_bitcoin/utils"
)
//区块链数据结构
type BlockChain struct {
	//所有的区块
	LastHash []byte
	DB *bolt.DB
}

// 生成区块链并持久化存储
// addresss 接受创世块奖励
func CreateBlockchainDB(address, nodeID string) *BlockChain {
	// 拼接数据库文件名字符串
	dbName := utils.DBName + "_" + nodeID + ".db"
	// 检查该文件是否已存在
	if dbExists(dbName){
		// 已存在就退出
		fmt.Println(dbName, " Blockchain already exists.")
		os.Exit(1)
	}
	// 创建创世交易
	genesisTx := NewCoinbaseTX(address, utils.GenesisCoinbaseData)
	// 创建创世区块
	// 初始化区块高度
	GenesisBlock := CreateGenesisBlcok(genesisTx, 0)
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
func (bc *BlockChain) AddNewBlock(txs []*Transaction)  {
	var err error
	//计算新区块数据
	preHash :=  bc.LastHash
	// 验证区块中的所有交易
	for _, tx := range txs {
		if ! bc.VerifyTx(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	//加入区块链
	err = bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BucketName))
		// 更新区块高度
		heightByte := bucket.Get([]byte(utils.BlockHeightKey))
		height := int64(utils.BytesToInt(heightByte) + 1)			// 高度+1
		err := bucket.Put([]byte(utils.BlockHeightKey), utils.Int64ToBytes(height))
		if err != nil {
			log.Panic(err)
		}
		// 生成新区块
		newBlock := NewBlock(txs, height, preHash)
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
		// 更新区块链记录的lastHash和高度
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
		fmt.Println(block)
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

// 创建区块链迭代器
func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		currentHash: bc.LastHash,
		db:          bc.DB,
	}
}


// FindUnspentTransactions 找到所有未花费输出的交易
// 遍历所有区块，每个区块的所有交易，每个交易的所有输出，如果输出的对象是本人那么就记录该交易
func (bc *BlockChain) FindUnspentTransactions(pubHashKey []byte) (txs []*Transaction) {
	// 标记数组，当遍历到单个交易时，其中的所有输入对应的上一个交易的输出都要排除掉
	// map: 交易Hash => output编号map(使用map而不用数组是为了不遍历)
	tapMap := make(map[string]map[int]bool, 0)
	bci := bc.NewIterator()
	block := bci.Next()
	// 遍历区块
	for block != nil {
		// 遍历区块中的交易
		for _, tx := range block.Transactions{
			// 遍历交易中的VOut,加入结果
			skip:
			for outputIdx, output := range tx.VOut {
				if _, has := tapMap[string(tx.TxHash)][outputIdx]; has{
					// 遍历其中的output标号数组，有的话直接结束
					continue skip
				}
				// 不存在的话检查是否可以解锁（是否为本人）
				if output.IsLockedWithKey(pubHashKey) {
					// 是，则添加到结果集
					txs = append(txs, tx)
				}
			}
			// 遍历Vin，打标记，优化遍历
			// 先排除coinbase交易，前面没有输出了，不需要标记
			if tx.IsCoinbase() {
				continue
			}
			for _, input := range tx.Vin {
				if outputMap, has := tapMap[string(input.OutputTxHash)]; has {
					outputMap[input.OutputIdx] = true
				}else{
					newMap := make(map[int]bool)
					newMap[input.OutputIdx] = true
					tapMap[string(input.OutputTxHash)] = newMap
				}
			}
		}
		block = bci.Next()
	}
	return
}

// 找到所有的UTXO
func (bc *BlockChain) FindUTXO(PubKeyHash []byte) (outputs []*TxOutput) {
	UtxoTxs := bc.FindUnspentTransactions(PubKeyHash)
	for _, tx := range UtxoTxs {
		for _, output := range tx.VOut {
			if output.IsLockedWithKey(PubKeyHash){
				outputs = append(outputs, output)
			}
		}
	}
	return
}

// 获取对应地址满足金额的UTXO集合
func (bc *BlockChain) FindSpendableOutputs(pubHashKey []byte, expectAmount int) (map[string][]int, int) {
	// 获取当前人的所有UTXO集合
	txs := bc.FindUnspentTransactions(pubHashKey)
	res := make(map[string][]int)
	cumulativePrice := 0
	out:
	for _, tx := range txs {
		for outputIdx, output := range tx.VOut {
			if output.IsLockedWithKey(pubHashKey) && cumulativePrice < expectAmount {
				cumulativePrice += output.Value
				res[string(tx.TxHash)] = append(res[string(tx.TxHash)], outputIdx)
				if cumulativePrice >= expectAmount {
					break out
				}
			}
		}

	}
	return res, cumulativePrice
}


/**
 * @Description: 查找一个交易根据交易Hash
 * @receiver bc
 * @param txHash
 * @return *Transaction
 * @return error
 */
func (bc *BlockChain) FindTransactionByTxHash(txHash []byte) (*Transaction, error) {
	iterator := bc.NewIterator()
	block := iterator.Next()
	for block != nil {
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.TxHash, txHash) == 0 {
				return tx, nil
			}
		}
		block = iterator.Next()
	}
	return nil, errors.New("Transaction Not Found!")
}

/**
 * @Description: 验证交易合法性,验证此交易中所有input前向TxHash的正确性
 * @receiver bc
 * @param tx
 * @return bool
 */
func (bc *BlockChain) VerifyTx(tx *Transaction) bool {
	preTxs := bc.GetPreTx(tx)
	// 验证单个交易合法性
	return tx.Verify(preTxs)
}

/**
 * @Description:  给交易签名
 * @receiver bc
 * @param tx
 * @param preKey  私钥
 */
func (bc *BlockChain) SignTransaction(tx *Transaction, priKey ecdsa.PrivateKey)  {
	preTxs := bc.GetPreTx(tx)
	// 签名单个交易
	tx.Sign(preTxs, priKey)
}

func (bc *BlockChain) GetPreTx(tx *Transaction) map[string]*Transaction {
	// 创建前向交易集的映射关系，方便签名
	preTxs := make(map[string]*Transaction)
	for _, input := range tx.Vin {
		// 查找
		tx, err := bc.FindTransactionByTxHash(input.OutputTxHash)
		if err != nil {
			log.Panic(err)
		}
		// 加入集合
		preTxs[hex.EncodeToString(tx.TxHash)] = tx
	}
	// 返回
	return preTxs
}