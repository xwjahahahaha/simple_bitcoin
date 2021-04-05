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
	BucketName string
	LastHashKey string
)


func init()  {
	file, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查配置文件！")
	}
	LoadBlockChainConfig(file)
	LoadConsensus(file)
	LoadDataBase(file)
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
	BucketName = file.Section("database").Key("BucketName").MustString("simpleBucket")
	LastHashKey = file.Section("database").Key("LastHashKey").MustString("LastHash")
}