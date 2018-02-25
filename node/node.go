package node

import (
	"os"
	"os/signal"
	"time"

	"github.com/frankh/nano/network"
	"github.com/frankh/nano/rpc"
	"github.com/frankh/nano/store"
	"github.com/frankh/nano/wallet"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	Net       *network.Network
	alarms    []*Alarm
	store     *store.Store
	rpc       *rpc.Server
	wallets   map[string]*wallet.Wallet
	walletsCh chan *wallet.Wallet
}

func NewNode(conf *store.Config) *Node {
	n := new(Node)

	n.Net = network.NewNetwork()
	n.alarms = make([]*Alarm, 1)
	n.store = store.NewStore(conf)
	n.wallets = make(map[string]*wallet.Wallet)
	n.walletsCh = make(chan *wallet.Wallet)

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

	n.alarms[0] = NewAlarm(AlarmFn(n.Net.SendKeepAlives), []interface{}{}, 20*time.Second)
	n.Net.ListenForUdp()
	n.rpc = rpc.NewServer(n.store, n.walletsCh)
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
			n.wallets[w.Id] = w
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
