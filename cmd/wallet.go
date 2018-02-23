package cmd

import (
	"fmt"

	"github.com/frankh/nano/wallet"

	"github.com/spf13/cobra"
)

var ()

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.AddCommand(walletCreateCmd)
}

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet management",
	Long:  `Create, modify and interact with wallets.`,
}

var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a wallet",
	Long:  `Generate a wallet seed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		w := wallet.NewWallet()
		s, err := w.GenerateSeed()
		if err != nil {
			return err
		}

		fmt.Println(s)

		return nil
	},
}
