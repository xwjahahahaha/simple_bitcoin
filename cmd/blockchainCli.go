package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"simple_bitcoin/blc"
)

func init() {
	rootCmd.AddCommand(CreateBlockChain, BlockChainPrint, AddBlock)
}



// 创建区块链
var	CreateBlockChain = &cobra.Command{
	Use: "createBC [nodeID]",
	Short: "create your blockchain",
	Long: "create your blockchain",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		blc.CreateBlockchainDB(args[0])
		fmt.Println("Done!")
	},

}

// 添加区块
var AddBlock = &cobra.Command{
	Use: "addBlock [nodeID] [data]",
	Short: "add your blockchain block",
	Long: "add your blockchain block",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取区块
		bc := blc.NewBlockchain(args[0])
		bc.AddNewBlock(args[1])
		fmt.Println("Done!")
		return nil
	},
}


// 输出区块链
var	BlockChainPrint = &cobra.Command{
	Use: "printBC [nodeID]",
	Short: "print your blockchain data",
	Long: "print your blockchain data",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bc := blc.NewBlockchain(args[0])
		bc.PrintBlockChain()
	},

}