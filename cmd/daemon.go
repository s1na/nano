package cmd

import (
	"log"
	"net"
	"path"

	"github.com/frankh/nano/node"
	"github.com/frankh/nano/store"

	"github.com/spf13/cobra"
)

var (
	InitialPeer string
	WorkDir     string
	TestNet     bool
)

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.Flags().StringVarP(&InitialPeer, "peer", "p", "::ffff:192.168.0.70", "Initial peer to make contact with")
	daemonCmd.Flags().StringVarP(&WorkDir, "work-dir", "d", "", "Directory to put generated files, e.g. db.")
	daemonCmd.Flags().BoolVarP(&TestNet, "testnet", "t", false, "Use test network configuration")
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Starts the node's daemon",
	Long:  `Starts a full Nano node as a long-running process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		n := node.NewNode()
		initialPeer := node.Peer{
			net.ParseIP(InitialPeer),
			7075,
		}
		n.Net.PeerList = []node.Peer{initialPeer}
		n.Net.PeerSet = map[string]bool{initialPeer.String(): true}

		if TestNet {
			log.Println("Using test network configuration")
			store.TestConfig.Path = path.Join(WorkDir, store.TestConfig.Path)
			store.Init(store.TestConfig)
		} else {
			store.LiveConfig.Path = path.Join(WorkDir, store.LiveConfig.Path)
			store.Init(store.LiveConfig)
		}

		n.Start()

		return nil
	},
}
