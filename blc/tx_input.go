package blc

type TxInput struct {
	// 与上一个交易相关的参数
	OutputTxHash []byte
	OutputIdx int
	// 签名与公钥
	Signature string
}

// 输入锁定方法
func (in *TxInput) CanUnlockedWith(unlockingData string) bool {
	// TODO 简单的锁定方法
	return in.Signature == unlockingData
}
