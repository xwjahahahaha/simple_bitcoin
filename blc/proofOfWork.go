package blc

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"simple_bitcoin/utils"
)

type ProofOfWork struct {
	Block *Block		// 当前要验证的区块
	Target *big.Int		// 区块的目标hash值(不能大于该值)
}

//创建pow对象
func NewProofOfWork(block *Block) *ProofOfWork {
	//1. 计算当前难度的最大值target，即满足当前难度要求的最大hash
	//例如： targetbit = 2 ， 那么 0100 0000 ... 0000 就是最大值，计算的结果要比这个小才满足
	// 创建末尾1
	oneBit := big.NewInt(1)
	// 移位的到最大值, 移位次数就是(HashDigits - targetBit)
	target := oneBit.Lsh(oneBit, utils.HashDigits - utils.TargetBit)

	return &ProofOfWork{
		Block: block,
		Target: target,
	}
}

//计算函数
//返回区块nonce以及Hash
func (pow *ProofOfWork) Run() (int64, []byte) {
	// 1. 拼接区块数据
	// 2. 根据nonce计算hash
	// 3. 比对hash， 满足退出，不满足修改nonce
	nonce := int64(0)						// 目标值
	hashBigInt := big.NewInt(nonce)			// 目标hash
	blockData := pow.PrepareData()
	for {
		hashBytes := sha256.Sum256(append(blockData, utils.Int64ToBytes(nonce)...))
		// 比对两个bigint的大小 x > y == 1 , x < y == -1 , x == y => 0
		if pow.Target.Cmp(hashBigInt.SetBytes(hashBytes[:])) == 1 {
			// 找到较小的hash， 满足了条件退出计算
			fmt.Printf("\r%x\n", hashBytes)
			break
		}
		fmt.Printf("\r%x", hashBytes)
		// 不满足，修改nonce继续
		nonce += 1
	}
	return nonce, hashBigInt.Bytes()
}

// 计算函数的区块数据统一化（统一为[]byte）函数
func (pow *ProofOfWork) PrepareData() []byte {
	return bytes.Join(
		[][]byte{
			utils.Int64ToBytes(pow.Block.BlockHeight),
		pow.Block.PreHash,
		pow.Block.Data,
			utils.Int64ToBytes(pow.Block.Timestamp),
		pow.Target.Bytes(),				// 当前区块难度
	},
		[]byte{},
	)
}


// PoW验证函数
// 验证目标区块是够合法（难度与区块的比对）
func (pow *ProofOfWork) IsValid() bool {
	// 转换区块Hash为big.Int
	blockHash := new(big.Int)
	blockHash.SetBytes(pow.Block.Hash)
	// 比较
	//	 -1 if x <  y
	//    0 if x == y
	//   +1 if x >  y
	if pow.Target.Cmp(blockHash) > 0 {
		return true
	}
	return false
}