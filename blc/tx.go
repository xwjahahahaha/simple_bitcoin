package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
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
	outputsDesc, acc :=  bc.FindSpendableOutputs(pubKeyHash, amount)

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
	// 签名交易
	bc.SignTransaction(tx, formWallet.PrivateKey)
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

/**
 * @Description: 签名整个交易
 * @receiver t
 * @param preTxs	相关的前置交易
 * @param priKey	私钥
 */
func (t *Transaction) Sign(preTxs map[string]*Transaction, priKey ecdsa.PrivateKey)  {
	// coinbase交易没有输入不签名
	if t.IsCoinbase() {
		return
	}
	// 验证此交易vin的所有前置交易都不为空
	for _, vin := range t.Vin {
		if preTxs[hex.EncodeToString(vin.OutputTxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	// 复制一个裁剪过的交易
	txCopy := t.TrimmedCopy()
	// 签名该交易所有的Vin
	// 默认一个交易的输入只有一个。。。
	for inputID, vin := range txCopy.Vin {
		// 获取当前vin的前置交易
		prevTx := preTxs[hex.EncodeToString(vin.OutputTxHash)]
		// 获取PubKeyHash并赋值到vin的Pubkey上（临时存储，为了取Hash）
		txCopy.Vin[inputID].Signature = nil
		txCopy.Vin[inputID].PubKey = prevTx.VOut[vin.OutputIdx].PubKeyHash
		// 交易取Hash
		// 这样就包含了出款人的pubHashKey（在input的pubkey中）和收款人的pubHashKey（在此交易的vout中）
		txCopy.SetHash()
		// 设置回nil
		txCopy.Vin[inputID].PubKey = nil

		// 对此交易整体进行签名(使用私钥)
		r, s, err := ecdsa.Sign(rand.Reader, &priKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		// 签名(也就是r和s的字节组合)  注意这里赋值给t，而不是txcopy
		t.Vin[inputID].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

/**
 * @Description: 验证整个交易
 * @receiver t
 * @param preTxs
 * @return bool
 */
func (t *Transaction) Verify(preTxs map[string]*Transaction) bool {
	if t.IsCoinbase() {
		return true
	}
	for _, vin := range t.Vin {
		if preTxs[hex.EncodeToString(vin.OutputTxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := t.TrimmedCopy()
	curve := elliptic.P256()

	// 验证交易中的每个vin
	for inputID, vin := range t.Vin {
		// 获取当前vin的前置交易
		prevTx := preTxs[hex.EncodeToString(vin.OutputTxHash)]
		// 获取PubKeyHash并赋值到vin的Pubkey上（临时存储，为了取Hash）
		txCopy.Vin[inputID].Signature = nil
		txCopy.Vin[inputID].PubKey = prevTx.VOut[vin.OutputIdx].PubKeyHash
		// 交易取Hash
		// 这样就包含了出款人的pubHashKey（在input的pubkey中）和收款人的pubHashKey（在此交易的vout中）
		txCopy.SetHash()
		// 设置回nil
		txCopy.Vin[inputID].PubKey = nil

		// 验证签名
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen)/2])	//r前半段
		s.SetBytes(vin.Signature[(sigLen)/2:])	//s后半段

		// 分割vin的公钥
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		// 验证
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}
	return true
}

/**
 * @Description:  修改要签名的整个交易其中的签名字段，为了签名时防止互相影响
 * @receiver t
 * @return *Transaction
 */
func (t *Transaction) TrimmedCopy() *Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput

	for _, vin := range t.Vin {
		// 注意就是将签名和公钥字段设置为nil
		inputs = append(inputs, &TxInput{vin.OutputTxHash, vin.OutputIdx, nil, nil})
	}

	for _, vout := range t.VOut {
		outputs = append(outputs, &TxOutput{vout.Value, vout.PubKeyHash})
	}
	// 拷贝
	txCopy := &Transaction{t.TxHash, inputs, outputs, t.Timestamp}

	return txCopy
}