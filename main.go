package main

import (
	"fmt"
	"simple_bitcoin/BLC"
)

func main() {

	genesisBlock := BLC.CreateGenesisBlcok("cqupt")
	fmt.Println(genesisBlock)
}
