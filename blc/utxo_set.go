package blc

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"simple_bitcoin/utils"
)

// UTXO未花费交易集合
type UTXOSet struct {
	Blockchain BlockChain
}

/**
 * @Description: 使用UTXO找到未花费输出，然后在数据库中进行存储。存储的形式：[TxHash] => utxos/outputs
 * @receiver u
 */
func (u *UTXOSet) Reindex()  {
	db := u.Blockchain.DB

	err := db.Update(func(tx *bolt.Tx) error {
		// 删除之前的bucket
		err := tx.DeleteBucket([]byte(utils.UtxoBucketName))
		// 创建新的bucket
		_, err = tx.CreateBucket([]byte(utils.UtxoBucketName))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})

	// 迭代区块链，获取所有的UTXO
	UTXOs := u.Blockchain.FindUTXO()

	// 存储到数据库
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.UtxoBucketName))
		// 存储每个key：交易 => value：outputs的映射序列
		for txHash, outs := range UTXOs {
			err = b.Put([]byte(txHash), outs.Serialize())
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

/**
 * @Description: 计算可以使用的金额以及对应的utxo
 * @receiver u
 * @param pubkeyHash	公钥Hash
 * @param amount
 * @return int
 * @return map[string][]int
 */
func (u *UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.DB
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.UtxoBucketName))
		c := b.Cursor()		// 创建游标遍历
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txHash := string(k)
			outs := DeSerializeOuts(v)
			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txHash] = append(unspentOutputs[txHash], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

/**
 * @Description: 获取当前人的所有utxo
 * @receiver u
 * @param pubkeyHash
 * @return []*TxOutput
 */
func (u *UTXOSet) FindUTXO(pubkeyHash []byte) []*TxOutput {
	var UTXO []*TxOutput
	db := u.Blockchain.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.UtxoBucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeSerializeOuts(v)
			for _, out := range outs.Outputs {
				// 判断是否为本人
				if out.IsLockedWithKey(pubkeyHash) {
					UTXO = append(UTXO, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return UTXO
}

/**
 * @Description: 当有新的区块生成时，根据区块更新UTXO存储，为了不重新遍历整条区块链
 * @receiver u
 * @param block
 */
func (u *UTXOSet) Update(block *Block)  {
	db := u.Blockchain.DB

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.UtxoBucketName))

		// 遍历新区块的交易
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {		// 注意，这里要确保当前的交易不是coinbase，如果是则只需要将其outputs全部添加即可
				// 遍历每个交易的输入
				for _, input := range tx.Vin {
					// 遍历交易的前置关联交易中的输出，做更新：被引用->删除，否则不变
					previousOutsBytes := b.Get(input.OutputTxHash)
					if previousOutsBytes == nil {
						log.Panic(fmt.Sprintf("Don't find %x in %s bucket", input.OutputTxHash, utils.UtxoBucketName))
					}
					previousOuts := DeSerializeOuts(previousOutsBytes)
					updateOuts := TxOutputs{}		// 创建新的output集合
					for outputIdx, output := range previousOuts.Outputs {
						// 以选择加入代替删除
						if input.OutputIdx != outputIdx {
							// 不等于，那么说明没有被新区块中的input引用，所以加入
							updateOuts.Outputs = append(updateOuts.Outputs, output)
						}
					}
					if len(updateOuts.Outputs) == 0 {
						// 都被删完了,直接在数据库中删除掉这个交易
						err := b.Delete(input.OutputTxHash)
						if err != nil {
							log.Panic(err)
						}
					}else {
						// 否则更新
						err := b.Put(input.OutputTxHash, updateOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}

			// 遍历每个交易的输出，直接添加到数据库
			outputs := TxOutputs{}
			for _, output := range tx.VOut {
				outputs.Outputs = append(outputs.Outputs, output)
			}
			err := b.Put(tx.TxHash, outputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}