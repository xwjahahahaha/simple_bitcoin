package main

import (
	"fmt"
	"simple_bitcoin/part1_basic_protoType/BLC"
)

func main()  {
	block := BLC.NewBlock([]byte("Genenis Block"), 1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	fmt.Println(block)
}
