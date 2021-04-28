package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"simple_bitcoin/utils"
	"strings"
	"time"
)

type Transaction struct {
	TxHash []byte
	Vin []*TxInput
	VOut []*TxOutput
	Timestamp int64
}

// NewUTXOTransaction creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == ""{
		data = utils.GenesisCoinbaseData
	}
	// 创建输入与输出
	txInput := &TxInput{
		OutputTxHash: INITIALLY_HASH,
		OutputIdx:    -1,				// 没有前置Hash
		Signature:    nil, 	// coinBase值
		PubKey: 		[]byte(data),
	}
	txOutput := &TxOutput{
		Value:        int(utils.CoinBaseReward),
		PubKeyHash: 	ResolveAddressToPubKeyHash(to),
	}

	// 创建genesis交易
	newTx :=  &Transaction{
		TxHash:    nil,
		Vin:       []*TxInput{txInput},
		VOut:      []*TxOutput{txOutput},
		Timestamp: time.Now().Unix(),
	}
	newTx.SetHash()
	return newTx
}

// 一般交易
func NewTransaction(from, to string,  amount int, bc *BlockChain) *Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput

	if len(from)==0 || len(to) == 0 || amount < 0 {
		log.Panic("参数错误")
		return nil
	}

	// 读取本地钱包
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	// 获取from对应的wallet
	formWallet := wallets.GetWallet(from)
	// 获取from的公钥Hash
	pubKeyHash := HashPubKey(formWallet.Publickey)

	// 获取当前能支付的金额
	outputsDesc, acc :=  bc.FindSpendableOutputs(pubKeyHash, amount)  // TODO FindSpendableOutputs

	if acc < amount {
		log.Panic("余额不足")
		return nil
	}

	// 创建交易
	// 创建Inputs
	for txHash, outputIdxAry := range outputsDesc {
		for _, outputIdx := range outputIdxAry {
			newInput := &TxInput{
				OutputTxHash: []byte(txHash),
				OutputIdx:    outputIdx,
				Signature:    nil,
				PubKey:		  formWallet.Publickey,			// 将form的公钥赋值
			}
			inputs = append(inputs, newInput)
		}
	}
	// 创建输出
	outputs = append(outputs, NewTxOutput(amount, to))
	// 找零(如果有的话)
	if acc > amount {
		outputs = append(outputs,NewTxOutput(acc - amount, from))
	}
	// 创建交易
	tx := &Transaction{
		TxHash:    nil,
		Vin:       inputs,
		VOut:      outputs,
		Timestamp: time.Now().Unix(),
	}
	tx.SetHash()
	return tx
}

// 交易序列化
func (t *Transaction) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(t)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// 反序列化
func Deserialize(Txbytes []byte) *Transaction {
	var t Transaction
	decoder := gob.NewDecoder(bytes.NewReader(Txbytes))
	err := decoder.Decode(&t)
	if err != nil {
		log.Panic(err)
	}
	return &t
}

// 计算交易Hash
func (t *Transaction) SetHash() {
	t.TxHash = []byte{}
	// 序列化
	bytes := t.Serialize()
	// sha256
	res := sha256.Sum256(bytes)
	t.TxHash = res[:]
}

// 是否为CoinBase
func (t *Transaction) IsCoinbase() bool {
	return len(t.Vin) == 1 && string(t.Vin[0].OutputTxHash) == string(INITIALLY_HASH) && t.Vin[0].OutputIdx == -1
}

// stdout打印
func (t *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("      txHash : %x\n", t.TxHash))
	lines = append(lines, fmt.Sprintf("      Vin : \n"))
	for i, input := range t.Vin {
		lines = append(lines, fmt.Sprintf("  [%d]:", i))
		lines = append(lines, fmt.Sprintln(input))
	}
	lines = append(lines, fmt.Sprintf("      VOut : \n"))
	for i, output := range t.VOut {
		lines = append(lines, fmt.Sprintf("  [%d]:", i))
		lines = append(lines, fmt.Sprintln(output))
	}
	lines = append(lines, fmt.Sprintf("      交易时间 : %s\n", time.Unix(t.Timestamp, 0).Format("2006-01-02 15:04:05")))
	return strings.Join(lines, "")
}

/**
 * @Description:	对当前交易签名
 * @receiver t
 * @param preTxs	input相关前缀交易集
 * @param priKey	签名者私钥
 */

// TODO
func (t *Transaction) Sign(preTxs map[string]*Transaction, priKey ecdsa.PrivateKey)  {

}

func (t *Transaction) Verify(preTxs map[string]*Transaction) bool {
	return false
}