package network

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/frankh/nano/store"

	log "github.com/sirupsen/logrus"
)

const packetSize = 512
const numberOfPeersToShare = 8

type Peer struct {
	IP   net.IP
	Port uint16
}

func (p *Peer) Addr() *net.UDPAddr {
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.IP.String(), p.Port))
	return addr
}

func (p *Peer) String() string {
	return fmt.Sprintf("%s:%d", p.IP.String(), p.Port)
}

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
			n.handleMessage(source, bytes.NewBuffer(buf[:c]))
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

func (n *Network) handleMessage(source string, buf *bytes.Buffer) {
	var header Header
	header.ReadHeader(bytes.NewBuffer(buf.Bytes()))
	if header.MagicNumber != MagicNumber {
		log.Printf("Ignored message. Wrong magic number %s", header.MagicNumber)
		return
	}

	sourcePeer := Peer{net.ParseIP(source), 7075}
	if !n.PeerSet[sourcePeer.String()] && source != n.LocalIP {
		n.PeerSet[sourcePeer.String()] = true
		n.PeerList = append(n.PeerList, sourcePeer)
		log.WithFields(log.Fields{
			"peer": sourcePeer.String(),
			"len":  len(n.PeerList),
		}).Info("Added new peer to list")
	}

	switch header.Type {
	case msgKeepalive:
		var m KeepAlive
		err := m.Read(buf)
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Warn("Failed to read keepalive")
		}

		err = m.Handle(n)
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Warn("Failed to handle keepalive")
		}
	case msgPublish:
		var m Publish
		err := m.Read(buf)
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Warn("Failed to read publish")
		} else {
			store.StoreBlock(m.ToBlock())
		}
	case msgConfirmAck:
		var m ConfirmAck
		err := m.Read(buf)
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Warn("Failed to read confirm")
		} else {
			store.StoreBlock(m.ToBlock())
		}
	default:
		log.WithFields(log.Fields{"type": header.Type}).Warn("Message type undefined, ignoring...")
	}
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
	buf := bytes.NewBuffer(nil)
	m.Write(buf)

	outConn.Write(buf.Bytes())
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
