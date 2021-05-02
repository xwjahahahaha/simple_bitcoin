package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"simple_bitcoin/blc"
	"strconv"
)

func init() {
	rootCmd.AddCommand(CreateBlockChain, BlockChainPrint, Send, GetBalance, CreateWallet)
}

// 创建区块链
var	CreateBlockChain = &cobra.Command{
	Use: "createBC [address] [nodeID]",
	Short: "create your blockchain",
	Long: "create your blockchain",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		address, nodeID := args[0], args[1]
		blc.CreateBlockchainDB(address, nodeID)
		fmt.Println("Blockchain ID : ", nodeID, " Done!")
	},
}

// 新增交易
var Send = &cobra.Command{
	Use: "send [from] [to] [amount] [nodeID]",
	Short: "add your blockchain block",
	Long: "add your blockchain block",
	Args: cobra.ExactArgs(4),

	Run: func(cmd *cobra.Command, args []string) {
		from, to, nodeID := args[0], args[1], args[3]
		amount, _ := strconv.Atoi(args[2])
		// 获取当前区块链数据
		bc := blc.NewBlockchain(nodeID)
		defer bc.DB.Close()
		// 创建交易
		newTX := blc.NewTransaction(from, to, amount, bc)
		// 创建区块,添加交易
		bc.AddNewBlock([]*blc.Transaction{newTX})
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
		defer bc.DB.Close()
		bc.PrintBlockChain()
	},
}

// 获取余额

var GetBalance = &cobra.Command{
	Use: "getBalance [address] [nodeID]",
	Short: "get your balance",
	Long: "get your balance",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取cli参数
		address, nodeID := args[0], args[1]
		// 验证地址
		if !blc.ValidateAddress(address) {
			return errors.New("Not a valid address")
		}
		// 加载区块链
		bc := blc.NewBlockchain(nodeID)
		defer bc.DB.Close()
		outputs := bc.FindUTXO(blc.ResolveAddressToPubKeyHash(address))
		amount := 0
		for _, output := range outputs {
			amount += output.Value
		}
		fmt.Println("----[", address, "]的余额为：", amount)
		return nil
	},
}

/**
 * @Description:  创建钱包
 */
var CreateWallet = &cobra.Command{
	Use: "createWallet",
	Short: "create a wallet",
	Long: "create a wallet",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error{
		wallets, _ := blc.NewWallets()
		// 创建公私钥对,返回对应的地址
		address := wallets.CreateWallet()
		// 保存到静态文件中
		wallets.SaveToFile()
		fmt.Println("Your new address : ", address)
		return nil
	},
}