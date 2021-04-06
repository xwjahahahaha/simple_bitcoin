package main

import "simple_bitcoin/cmd"
func main() {

	//BlockChain := blc.CreateBlockchainDB("xwj")
	////
	////
	//defer BlockChain.DB.Close()

	//bc := blc.NewBlockchain("x1wj")
	//fmt.Printf("%x\n", bc.LastHash)



	//添加新区块
	//BlockChain.AddNewBlock("a send $1 to b")
	//BlockChain.AddNewBlock("a send $2 to b")
	//BlockChain.AddNewBlock("a send $3 to b")
	//fmt.Printf("本区块Hash值: %x\n", BlockChain.LastHash)

	//BlockChain.PrintBlockChain()
	cmd.Execute()
}
