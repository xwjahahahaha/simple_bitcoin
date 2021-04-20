package blc

type TxInput struct {
	// 与上一个交易相关的参数
	OutputTxHash []byte
	OutputIdx int
	// 签名与公钥
	Signature []byte
	PubKey []byte
}

//func (txInput *TxInput) String() string {
//
//}
