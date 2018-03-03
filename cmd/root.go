package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/frankh/nano/config"

	"github.com/spf13/cobra"
)

var (
	DataDir string
	TestNet bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&DataDir, "data-dir", "d", "", "Directory to put generated files, e.g. db.")
	rootCmd.PersistentFlags().BoolVarP(&TestNet, "testnet", "t", false, "Use test network configuration")
}

var rootCmd = &cobra.Command{
	Use:   "nanode",
	Short: "Nanode is a Go-based Nano node",
	Long:  `Nanode is a Go-based Nano node`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if TestNet {
			DataDir = path.Join(DataDir, config.TestDBName)
		} else {
			DataDir = path.Join(DataDir, config.DBName)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
