package cmd

import (
	"net"

	"github.com/s1na/nano/config"
	"github.com/s1na/nano/network"
	"github.com/s1na/nano/node"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	InitialPeer string
	Verbose     bool
)

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.Flags().StringVarP(&InitialPeer, "peer", "p", "::ffff:192.168.0.70", "Initial peer to make contact with")
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

		conf := &config.Config{
			DataDir: DataDir,
		}
		if TestNet {
			log.Info("Using test network configuration")
			conf.TestNet = true
		}

		n := node.NewNode(conf)
		initialPeer := network.Peer{
			net.ParseIP(InitialPeer),
			7075,
		}
		n.Net.PeerList = []network.Peer{initialPeer}
		n.Net.PeerSet = map[string]bool{initialPeer.String(): true}

		n.Start()

		return nil
	},
}
