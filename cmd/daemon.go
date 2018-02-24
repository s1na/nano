package cmd

import (
	"net"
	"path"

	"github.com/frankh/nano/node"
	"github.com/frankh/nano/store"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	InitialPeer string
	WorkDir     string
	TestNet     bool
	Verbose     bool
)

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.Flags().StringVarP(&InitialPeer, "peer", "p", "::ffff:192.168.0.70", "Initial peer to make contact with")
	daemonCmd.Flags().StringVarP(&WorkDir, "work-dir", "d", "", "Directory to put generated files, e.g. db.")
	daemonCmd.Flags().BoolVarP(&TestNet, "testnet", "t", false, "Use test network configuration")
	daemonCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose mode")
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Starts the node's daemon",
	Long:  `Starts a full Nano node as a long-running process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		var conf store.Config
		if TestNet {
			log.Info("Using test network configuration")
			store.TestConfig.Path = path.Join(WorkDir, store.TestConfig.Path)
			conf = store.TestConfig
		} else {
			store.LiveConfig.Path = path.Join(WorkDir, store.LiveConfig.Path)
			conf = store.LiveConfig
		}

		n := node.NewNode(&conf)
		initialPeer := node.Peer{
			net.ParseIP(InitialPeer),
			7075,
		}
		n.Net.PeerList = []node.Peer{initialPeer}
		n.Net.PeerSet = map[string]bool{initialPeer.String(): true}

		n.Start()

		return nil
	},
}
