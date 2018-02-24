package cmd

import (
	"fmt"

	"github.com/frankh/nano/store"
	"github.com/frankh/nano/wallet"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var ()

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountCreateCmd)
	accountCmd.AddCommand(accountGetCmd)
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account management",
	Long:  `Create, modify and interact with accounts.`,
}

var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an account in a wallet",
	Long:  `Create an account in a wallet, and display its address.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wid := args[0]

		store := store.NewStore(&store.TestConfig)
		if err := store.Start(); err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		w, err := ws.GetWallet(wid)
		if err != nil {
			return err
		}

		a := w.NewAccount()
		if err = ws.SetWallet(w); err != nil {
			return err
		}

		fmt.Println(a.Address())

		return nil
	},
}

var accountGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch information about an account",
	Long:  `Display information stored about an account.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := args[0]

		store := store.NewStore(&store.TestConfig)
		if err := store.Start(); err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		a, err := ws.GetAccount(addr)
		if err != nil {
			return err
		}

		if a == nil {
			return errors.New("Account not found in any wallet")
		}

		fmt.Println(a.String())

		return nil
	},
}
