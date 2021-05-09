package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	GenesisCoinbaseData string

	TargetBit uint		// 难度
	HashDigits uint

	DBName string
	BlockBucketName string
	UtxoBucketName string
	LastHashKey string
	BlockHeightKey string

	RootCmd string
	RootShort string
	RootLong string

	CoinBaseReward uint
	WalletsFileName string


	PubKeyVersion int
	AddressCheckSumLen int  // 截取校验和Hash后字节数


)


func init()  {
	file, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查配置文件！")
	}
	LoadBlockChainConfig(file)
	LoadConsensus(file)
	LoadDataBase(file)
	LoadCmdConfig(file)
	LoadTransaction(file)
	loadEncrypto(file)
}

func LoadBlockChainConfig(file *ini.File)  {
	GenesisCoinbaseData = file.Section("genesis").Key("GenesisCoinbaseData").MustString("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks")
}

func LoadConsensus(file *ini.File)  {
	TargetBit = file.Section("pow").Key("TargetBit").MustUint(16)
	HashDigits = file.Section("pow").Key("HashDigits").MustUint(256)
}

func LoadDataBase(file *ini.File)  {
	DBName = file.Section("database").Key("DBName").MustString("simpleCoin")
	BlockBucketName = file.Section("database").Key("BlockBucketName").MustString("simpleBucket")
	UtxoBucketName = file.Section("database").Key("UtxoBucketName").MustString("chainState")
	LastHashKey = file.Section("database").Key("LastHashKey").MustString("LastHash")
	BlockHeightKey = file.Section("database").Key("BlockHeightKey").MustString("BlockHeight")
}

func LoadCmdConfig(file *ini.File)  {
	RootCmd = file.Section("cmd").Key("RootCmd").MustString("btc")
	RootShort = file.Section("cmd").Key("RootShort").MustString("Simple bitcoin by gump")
	RootLong = file.Section("cmd").Key("RootLong").MustString("Simple bitcoin by gump")
}

func LoadTransaction(file *ini.File)  {
	CoinBaseReward = file.Section("tx").Key("CoinBaseReward").MustUint(10)
	WalletsFileName = file.Section("tx").Key("WalletsFileName").MustString("wallets.dat")
}

func loadEncrypto(file *ini.File)  {
	PubKeyVersion = file.Section("encrypto").Key("PubKeyVersion").MustInt(1)
	AddressCheckSumLen = file.Section("encrypto").Key("AddressCheckSumLen").MustInt(4)
}