package main

import (
	"fmt"
	"simple_bitcoin/BLC"
)

func main() {

	blockChain := BLC.CreateNewBC()
	//fmt.Println(newBlockChain)
	//fmt.Println(newBlockChain.Blocks)
	//fmt.Println(newBlockChain.Blocks[0])


	//添加新区块
	blockChain.AddNewBlock("a send $1 to b")

	fmt.Printf("%x\n", blockChain.Blocks[1].Hash)
	block := blockChain.Blocks[1]
	bytes := block.Serialization()
	fmt.Printf("序列化之后的结果：%v\n", bytes)
	block = BLC.Deserialization(bytes)
	fmt.Printf("%x\n", block.Hash)

	//fmt.Println(blockChain.Blocks[1].Nonce)
	//pow := BLC.NewProofOfWork(blockChain.Blocks[1])
	//fmt.Println("target:", pow.Target)
	//fmt.Println(pow.IsValid())
	//blockChain.AddNewBlock("a send $2 to b")
	//blockChain.AddNewBlock("a send $3 to b")
	//fmt.Println(blockChain)
	//fmt.Println(blockChain.Blocks)
	//fmt.Println(blockChain.Blocks[2])




}
