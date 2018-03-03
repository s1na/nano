package cmd

import (
	"errors"
	"fmt"

	"github.com/s1na/nano/store"
	"github.com/s1na/nano/types"
	"github.com/s1na/nano/wallet"

	"github.com/spf13/cobra"
)

var ()

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.AddCommand(walletCreateCmd)
	//walletCmd.AddCommand(walletImportCmd)
	walletCmd.AddCommand(walletGetCmd)
	walletCmd.AddCommand(walletListCmd)
	walletCmd.AddCommand(walletAddCmd)
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

		store := store.NewStore(DataDir)
		if err = store.Start(); err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		if err = ws.SetWallet(w); err != nil {
			return err
		}

		fmt.Println(id.Hex())

		return nil
	},
}

/*
var walletImportCmd = &cobra.Command{
	Use:   "import WALLET_ID WALLET_SEED",
	Short: "Import a wallet",
	Long:  `Import a wallet via its id and seed.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr := args[0]
		seedStr := args[1]

		w := wallet.NewWallet()
		w.Seed = seed
		w.Id = id

		store := store.NewStore(DataDir)
		if err := store.Start(); err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		if err := ws.SetWallet(w); err != nil {
			return err
		}

		fmt.Printf("Stored %s\n", w.Id.Hex())

		return nil
	},
}*/

var walletGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch information about a wallet",
	Long:  `Display information stored about a wallet.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := types.PubKeyFromHex(args[0])
		if err != nil {
			return err
		}

		store := store.NewStore(DataDir)
		err = store.Start()
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
		store := store.NewStore(DataDir)
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

var walletAddCmd = &cobra.Command{
	Use:   "add WALLET_ID PRIVATE_KEY",
	Short: "Add an adhoc account",
	Long:  `Add an adhoc private key to wallet, and display its address.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		wid, err := types.PubKeyFromHex(args[0])
		if err != nil {
			return err
		}

		key, err := types.PrvKeyFromString(args[1])
		if err != nil {
			return err
		}

		store := store.NewStore(DataDir)
		err = store.Start()
		if err != nil {
			return err
		}

		ws := wallet.NewWalletStore(store)
		wal, err := ws.GetWallet(wid)
		if err != nil {
			return err
		}

		pub, prv, err := types.KeypairFromPrvKey(key)
		if err != nil {
			return err
		}

		wal.InsertAdhoc(pub, prv)
		if err = ws.SetWallet(wal); err != nil {
			return errors.New("internal error")
		}

		fmt.Println(pub.Address())

		return nil
	},
}
