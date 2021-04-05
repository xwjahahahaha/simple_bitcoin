package main

import (
	"simple_bitcoin/BLC"
)

func main() {

	blockChain := BLC.CreateNewBC()
	defer blockChain.DB.Close()
	//fmt.Println(newBlockChain)
	//fmt.Println(newBlockChain.Blocks)
	//fmt.Println(newBlockChain.Blocks[0])


	//添加新区块
	blockChain.AddNewBlock("a send $1 to b")
	blockChain.AddNewBlock("a send $2 to b")
	blockChain.AddNewBlock("a send $3 to b")

	blockChain.PrintBlockChain()
}
