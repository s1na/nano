package node

import (
	"bytes"
	"log"
	"math/rand"
	"net"

	"github.com/frankh/nano/store"
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
	n.stop = make(chan bool)

	return n
}

func (n *Network) Stop() {
	n.stop <- true
}

func (n *Network) ListenForUdp() {
	go n.listenForUdp()
}

func (n *Network) listenForUdp() {
	log.Printf("Listening for udp packets on 7075")
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
			log.Println("UPDListener got stop message, shutting down...")
			break
		default:
		}
	}
}

func (n *Network) handleMessage(source string, buf *bytes.Buffer) {
	var header MessageHeader
	header.ReadHeader(bytes.NewBuffer(buf.Bytes()))
	if header.MagicNumber != MagicNumber {
		log.Printf("Ignored message. Wrong magic number %s", header.MagicNumber)
		return
	}

	sourcePeer := Peer{net.ParseIP(source), 7075}
	if !n.PeerSet[sourcePeer.String()] && source != n.LocalIP {
		n.PeerSet[sourcePeer.String()] = true
		n.PeerList = append(n.PeerList, sourcePeer)
		log.Printf("Added new peer to list: %s, now %d peers", sourcePeer.String(), len(n.PeerList))
	}

	switch header.MessageType {
	case Message_keepalive:
		var m MessageKeepAlive
		err := m.Read(buf)
		if err != nil {
			log.Printf("Failed to read keepalive: %s", err)
		}
		log.Printf("Read keepalive from %s", source)
		err = m.Handle(n)
		if err != nil {
			log.Printf("Failed to handle keepalive")
		}
	case Message_publish:
		var m MessagePublish
		err := m.Read(buf)
		if err != nil {
			log.Printf("Failed to read publish: %s", err)
		} else {
			store.StoreBlock(m.ToBlock())
		}
	case Message_confirm_ack:
		var m MessageConfirmAck
		err := m.Read(buf)
		if err != nil {
			log.Printf("Failed to read confirm: %s", err)
		} else {
			store.StoreBlock(m.ToBlock())
		}
	default:
		log.Printf("Ignored message. Cannot handle message type %d\n", header.MessageType)
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

	m := CreateKeepAlive(randomPeers)
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
