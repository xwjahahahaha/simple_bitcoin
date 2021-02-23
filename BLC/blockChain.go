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

//区块链中添加新的区块
func (bc *BlockChain) AddNewBlock(data string)  {
	length := len(bc.Blocks)
	lastBlock := bc.Blocks[length-1]
	//计算新区块数据
	newHeight, preHash := int64(length), lastBlock.Hash
	//生成新区块
	newBlock := NewBlock(data, newHeight, preHash)
	//加入区块链
	bc.Blocks = append(bc.Blocks, newBlock)
}

