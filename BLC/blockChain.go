package BLC

//区块链数据结构
type BlockChain struct {
	//所有的区块
	Blocks []*Block

}

//创建新区块链
func CreateNewBC() *BlockChain {
	//创建创世区块
	GenesisBlock := CreateGenesisBlcok("xwj")
	//添加到区块链中
	return &BlockChain{
		[]*Block{GenesisBlock},
	}
}

