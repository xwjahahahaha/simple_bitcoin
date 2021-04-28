package blc

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"simple_bitcoin/utils"
)

/**
 * @Description: 多钱包，地址 => 单钱包
 */
type Wallets struct {
	Wallets map[string]*Wallet
}

/**
 * @Description: 读取本地钱包
 * @return *Wallets
 * @return error
 */
func NewWallets()  (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	// 加载本地文件
	err := wallets.LoadFromFile()

	return &wallets, err
}

/**
 * @Description: 在钱包中创建一对公私钥
 * @receiver ws
 * @return string 地址
 */
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()
	ws.Wallets[string(address)] = wallet
	return string(address)
}

/**
 * @Description: 加载本地文件 loads wallets from file
 * @receiver ws
 * @return error
 */
func (ws *Wallets) LoadFromFile() error {
	// 这里存储的文件没有使用blot，就是序列化存储为一个dat文件
	// 1. 检查文件是否存在
	if _, err := os.Stat(utils.WalletsFileName); os.IsNotExist(err) {
		return err
	}
	// 2. 读取文件
	fileBytes, err := ioutil.ReadFile(utils.WalletsFileName)
	if err != nil {
		log.Panic(err)
	}
	// 3. 解码/反序列化
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileBytes))
	// 4. 加载到wallets结构体
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	// 5. 赋值
	ws.Wallets = wallets.Wallets
	return nil
}

/**
 * @Description: 保存到文件
 * @receiver ws
 */
func (ws *Wallets) SaveToFile()  {
	// 1. 序列化
	content := new(bytes.Buffer)       // 创建缓冲区
	gob.Register(elliptic.P256())      // 注册
	encoder := gob.NewEncoder(content) // 创建encoder
	err := encoder.Encode(ws)          // 序列化wallets
	if err != nil {
		log.Panic(err)
	}
	// 2. 打开文件,写入
	err = ioutil.WriteFile(utils.WalletsFileName, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

/**
 * @Description:  根据地址获取对应的Wallet
 * @receiver ws
 * @param address
 * @return *Wallet
 */
func (ws *Wallets) GetWallet(address string) *Wallet {
	return ws.Wallets[address]
}

/**
 * @Description: 获取钱包的所有地址
 * @receiver ws
 * @return []string
 */
func (ws *Wallets) GetAddresses() (addresses []string) {
	for k, _ := range ws.Wallets {
		addresses = append(addresses, k)
	}
	return
}

