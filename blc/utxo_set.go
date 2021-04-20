package blc

// UTXO未花费交易集合
type UTXOSet struct {
	Blockchain BlockChain
}

// 找到所有为花费的交易输出
//func (u *UTXOSet) FindUnSpendableOutputs() (int, map[string][]int) {
//
//}