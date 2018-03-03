package cmd

import (
	"fmt"

	"github.com/frankh/nano/account"
	"github.com/frankh/nano/store"
	"github.com/frankh/nano/types"
	"github.com/frankh/nano/wallet"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

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
		wid, err := types.PubKeyFromHex(args[0])
		if err != nil {
			return err
		}

		store := store.NewStore(DataDir)
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
		pub, err := types.PubKeyFromAddress(args[0])
		if err != nil {
			return err
		}

		store := store.NewStore(DataDir)
		if err := store.Start(); err != nil {
			return err
		}

		as := account.NewAccountStore(store)
		a, err := as.GetAccount(pub)
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
