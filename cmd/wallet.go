package cmd

import (
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
	walletCmd.AddCommand(walletListCmd)
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
		if err = store.Start(); err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		if err = ws.SetWallet(w); err != nil {
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

		fmt.Println(w.String())

		return nil
	},
}

var walletListCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallets",
	Long:  `Display information about all locally stored wallets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store := store.NewStore(&store.TestConfig)
		err := store.Start()
		if err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		wallets, err := ws.GetWallets()
		if err != nil {
			return err
		}

		for _, w := range wallets {
			fmt.Println(w.String())
		}

		return nil
	},
}
