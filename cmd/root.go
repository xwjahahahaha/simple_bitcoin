package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simple_bitcoin/utils"
)

var rootCmd = &cobra.Command{
	Use: utils.RootCmd,
	Short: utils.RootShort,
	Long: utils.RootLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bitcoin stable 1.0")
	},
}

func Execute()  {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}