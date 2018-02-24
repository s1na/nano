package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/frankh/nano/store"
	"github.com/frankh/nano/wallet"

	"github.com/spf13/cobra"
)

var ()

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.AddCommand(walletCreateCmd)
	walletCmd.AddCommand(walletGetCmd)
}

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet management",
	Long:  `Create, modify and interact with wallets.`,
}

var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a wallet",
	Long:  `Create a wallet, and display its ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		w := wallet.NewWallet()
		id, err := w.Init()
		if err != nil {
			return err
		}

		store := store.NewStore(&store.TestConfig)
		err = store.Start()
		if err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		err = ws.SetWallet(w)
		if err != nil {
			return err
		}

		fmt.Println(id)

		return nil
	},
}

var walletGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch information about a wallet",
	Long:  `Display information stored about a wallet.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		store := store.NewStore(&store.TestConfig)
		err := store.Start()
		if err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		w, err := ws.GetWallet(id)
		if err != nil {
			return err
		}

		output, err := json.Marshal(w)
		if err != nil {
			return err
		}

		fmt.Println(string(output))

		return nil
	},
}
