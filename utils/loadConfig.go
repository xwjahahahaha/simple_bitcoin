package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	GenesisString string

	TargetBit uint		// 难度
	HashDigits uint

)

func init()  {
	file, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查配置文件！")
	}
	LoadBlockChainConfig(file)
	LoadConsensus(file)
}

func LoadBlockChainConfig(file *ini.File)  {
	GenesisString = file.Section("genesis").Key("GenesisString").MustString("genesis")

}

func LoadConsensus(file *ini.File)  {
	TargetBit = file.Section("pow").Key("TargetBit").MustUint(16)
	HashDigits = file.Section("pow").Key("HashDigits").MustUint(256)
}