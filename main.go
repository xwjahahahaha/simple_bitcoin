package main

import (
	"fmt"
	"simple_bitcoin/BLC"
)

func main() {

	newBlockChain := BLC.CreateNewBC()
	fmt.Println(newBlockChain)
	fmt.Println(newBlockChain.Blocks)
	fmt.Println(newBlockChain.Blocks[0])
}
