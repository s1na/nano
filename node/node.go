package node

import (
	"os"
	"os/signal"
	"time"

	"github.com/s1na/nano/blocks"
	"github.com/s1na/nano/config"
	"github.com/s1na/nano/ledger"
	"github.com/s1na/nano/network"
	"github.com/s1na/nano/rpc"
	"github.com/s1na/nano/store"
	"github.com/s1na/nano/wallet"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	Net       *network.Network
	alarms    []*Alarm
	store     *store.Store
	rpc       *rpc.Server
	ledger    *ledger.Ledger
	wallets   map[string]*wallet.Wallet
	walletsCh chan *wallet.Wallet
	blocksCh  chan blocks.Block
}

func NewNode(conf *config.Config) *Node {
	n := new(Node)

	config.Conf = conf
	blocks.GenesisBlock = blocks.LiveGenesisBlock
	if conf.TestNet {
		blocks.GenesisBlock = blocks.TestGenesisBlock
	}

	n.Net = network.NewNetwork()
	n.store = store.NewStore(conf.DataDir)
	n.ledger = ledger.NewLedger(n.store)
	n.alarms = make([]*Alarm, 1)
	n.wallets = make(map[string]*wallet.Wallet)
	n.walletsCh = make(chan *wallet.Wallet)
	n.blocksCh = make(chan blocks.Block)

	return n
}

func (n *Node) Start() {
	if err := n.store.Start(); err != nil {
		log.Fatal(err)
	}

	// Fetch previously stored data in db
	if err := n.syncFromStore(); err != nil {
		log.Fatal(err)
	}

	if err := n.ledger.Init(); err != nil {
		log.Fatal(err)
	}

	n.alarms[0] = NewAlarm(AlarmFn(n.Net.SendKeepAlives), []interface{}{}, 20*time.Second)
	n.Net.ListenForUdp()
	n.rpc = rpc.NewServer(n.store, n.walletsCh, n.blocksCh)
	n.rpc.Start()

	n.loop()
}

func (n *Node) Stop() {
	n.rpc.Stop()
	n.alarms[0].Stop()
	n.Net.Stop()
}

func (n *Node) loop() {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, os.Kill)

	log.Info("Starting node loop")

	for {
		select {
		case s := <-sigCh:
			log.WithFields(log.Fields{"signal": s.String()}).Info("Caught signal, shutting down...")
			n.Stop()
			return
		case w := <-n.walletsCh:
			log.WithFields(log.Fields{"wallet": w.Id}).Info("Adding wallet to node")
			n.wallets[w.Id.Hex()] = w
		case b := <-n.blocksCh:
			err := n.ledger.AddBlock(b)
			if err != nil {
				log.WithFields(log.Fields{"block": b.Hash(), "err": err.Error()}).Warn("Failed adding block to ledger")
			}
		}
	}

	log.Info("Stopping node loop")
}

func (n *Node) syncFromStore() error {
	ws := wallet.NewWalletStore(n.store)
	wallets, err := ws.GetWallets()
	if err != nil {
		return err
	}
	n.wallets = wallets

	return nil
}
