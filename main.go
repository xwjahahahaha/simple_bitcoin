package main

import (
	"fmt"
	"simple_bitcoin/BLC"
)

func main()  {

	block := BLC.NewBlock([]byte("Genenis Block"), 1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	block.SetBlockHash()
	fmt.Println(block)
}
