package network

import (
	"math/rand"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const packetSize = 512
const numberOfPeersToShare = 8

type Network struct {
	PeerList []Peer
	PeerSet  map[string]bool
	LocalIP  string
	stop     chan bool
}

func NewNetwork() *Network {
	n := new(Network)

	n.PeerList = make([]Peer, 0, 5)
	n.PeerSet = make(map[string]bool)
	n.LocalIP = getOutboundIP().String()
	n.stop = make(chan bool, 1)

	return n
}

func (n *Network) Stop() {
	n.stop <- true
	time.Sleep(500 * time.Millisecond)
}

func (n *Network) ListenForUdp() {
	go n.listenForUdp()
}

func (n *Network) listenForUdp() {
	log.Info("Listening for udp packets on 7075")
	ln, err := net.ListenPacket("udp", ":7075")
	if err != nil {
		panic(err)
	}

	buf := make([]byte, packetSize)

	for {
		c, addr, err := ln.ReadFrom(buf)
		if err != nil {
			continue
		}

		source := addr.(*net.UDPAddr).IP.String()
		if c > 0 {
			n.handleMessage(source, buf[:c])
		}

		// Check whether to stop, without blocking
		// TODO: Figure out a way to overcome ReadFrom blocking
		select {
		case _ = <-n.stop:
			log.Info("UPDListener got stop message, shutting down...")
			break
		default:
		}
	}
}

func (n *Network) AddPeer(p Peer) {
	if !n.PeerSet[p.String()] && p.IP.String() != n.LocalIP {
		n.PeerSet[p.String()] = true
		n.PeerList = append(n.PeerList, p)
		log.WithFields(log.Fields{
			"peer": p.String(),
			"len":  len(n.PeerList),
		}).Info("Added new peer to list")
	}
}

func (n *Network) handleMessage(source string, data []byte) {
	msg := new(Message)
	if err := msg.Unmarshal(data); err != nil {
		log.WithFields(log.Fields{"source": source, "err": err.Error()}).Warn("Failed to unmarshal message")
		return
	}

	sp := Peer{net.ParseIP(source), 7075}
	n.AddPeer(sp)

	switch m := msg.Body.(type) {
	case *KeepAlive:
		for _, peer := range m.Peers {
			n.AddPeer(peer)
		}
	case *Publish:
		//store.StoreBlock(m.ToBlock())
	case *ConfirmAck:
		//store.StoreBlock(m.ToBlock())
	}

	return
}

func (n *Network) SendKeepAlive(peer Peer) error {
	addr := peer.Addr()
	randomPeers := make([]Peer, 0, 2)

	outConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}

	randIndices := rand.Perm(len(n.PeerList))
	for j, i := range randIndices {
		if j == numberOfPeersToShare {
			break
		}
		randomPeers = append(randomPeers, n.PeerList[i])
	}

	m := NewKeepAlive(randomPeers)
	msg := NewMessage(msgKeepalive, m)
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	outConn.Write(data)

	return nil
}

func (n *Network) SendKeepAlives(params []interface{}) {
	for _, peer := range n.PeerList {
		// TODO: Handle errors
		n.SendKeepAlive(peer)
	}
}

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
