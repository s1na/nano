package network

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	MagicNumber = [2]byte{'R', 'C'}
)

const (
	VersionMax   = 0x06
	VersionUsing = 0x06
	VersionMin   = 0x06
)

const (
	msgInvalid byte = iota
	msgNotAType
	msgKeepalive
	msgPublish
	msgConfirmReq
	msgConfirmAck
	msgBulkPull
	msgBulkPush
	msgFrontierReq
)

const (
	blockInvalid byte = iota
	blockNotABlock
	blockSend
	blockReceive
	blockOpen
	blockChange
)

type Header struct {
	MagicNumber  [2]byte
	VersionMax   byte
	VersionUsing byte
	VersionMin   byte
	Type         byte
	Extensions   byte
	BlockType    byte
}

func NewHeader(t byte) *Header {
	h := new(Header)

	h.MagicNumber = MagicNumber
	h.VersionMax = VersionMax
	h.VersionUsing = VersionUsing
	h.VersionMin = VersionMin
	h.Type = t

	return h
}

func (m *Header) WriteHeader(buf *bytes.Buffer) error {
	var errs []error
	errs = append(errs,
		buf.WriteByte(m.MagicNumber[0]),
		buf.WriteByte(m.MagicNumber[1]),
		buf.WriteByte(m.VersionMax),
		buf.WriteByte(m.VersionUsing),
		buf.WriteByte(m.VersionMin),
		buf.WriteByte(m.Type),
		buf.WriteByte(m.Extensions),
		buf.WriteByte(m.BlockType),
	)

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Header) ReadHeader(buf *bytes.Buffer) error {
	var errs []error
	var err error
	// I really hate go error handling sometimes
	m.MagicNumber[0], err = buf.ReadByte()
	errs = append(errs, err)
	m.MagicNumber[1], err = buf.ReadByte()
	errs = append(errs, err)
	m.VersionMax, err = buf.ReadByte()
	errs = append(errs, err)
	m.VersionUsing, err = buf.ReadByte()
	errs = append(errs, err)
	m.VersionMin, err = buf.ReadByte()
	errs = append(errs, err)
	m.Type, err = buf.ReadByte()
	errs = append(errs, err)
	m.Extensions, err = buf.ReadByte()
	errs = append(errs, err)
	m.BlockType, err = buf.ReadByte()
	errs = append(errs, err)

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

type KeepAlive struct {
	Header
	Peers []Peer
}

func NewKeepAlive(peers []Peer) *KeepAlive {
	m := new(KeepAlive)

	m.Header = *NewHeader(msgKeepalive)
	m.Peers = peers

	return m
}

func (m *KeepAlive) Handle(network *Network) error {
	for _, peer := range m.Peers {
		if peer.IP.String() == network.LocalIP {
			continue
		}

		if !network.PeerSet[peer.String()] {
			network.PeerSet[peer.String()] = true
			network.PeerList = append(network.PeerList, peer)
			log.WithFields(log.Fields{
				"peer": peer.String(),
				"len":  len(network.PeerList),
			}).Info("Added new peer to list")
		}
	}

	return nil
}

func (m *KeepAlive) Read(buf *bytes.Buffer) error {
	var header Header
	err := header.ReadHeader(buf)
	if err != nil {
		return err
	}

	if header.Type != msgKeepalive {
		return errors.New("Tried to read wrong message type")
	}

	m.Header = header
	m.Peers = make([]Peer, 0)

	for {
		peerPort := make([]byte, 2)
		peerIp := make(net.IP, net.IPv6len)
		n, err := buf.Read(peerIp)
		if n == 0 {
			break
		}
		if err != nil {
			return err
		}
		n2, err := buf.Read(peerPort)
		if err != nil {
			return err
		}
		if n < net.IPv6len || n2 < 2 {
			return errors.New("Not enough ip bytes")
		}

		m.Peers = append(m.Peers, Peer{peerIp, binary.LittleEndian.Uint16(peerPort)})
	}

	return nil
}

func (m *KeepAlive) Write(buf *bytes.Buffer) error {
	err := m.Header.WriteHeader(buf)
	if err != nil {
		return err
	}

	for _, peer := range m.Peers {
		_, err = buf.Write(peer.IP)
		if err != nil {
			return err
		}
		portBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(portBytes, peer.Port)
		if err != nil {
			return err
		}
		_, err = buf.Write(portBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

type ConfirmAck struct {
	Header
	Vote
}

func (m *ConfirmAck) Read(buf *bytes.Buffer) error {
	err := m.Header.ReadHeader(buf)
	if err != nil {
		return err
	}

	if m.Header.Type != msgConfirmAck {
		return errors.New("Tried to read wrong message type")
	}
	err = m.Vote.Read(m.Header.BlockType, buf)
	if err != nil {
		return err
	}

	return nil
}

func (m *ConfirmAck) Write(buf *bytes.Buffer) error {
	err := m.Header.WriteHeader(buf)
	if err != nil {
		return err
	}

	err = m.Vote.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

type ConfirmReq struct {
	Header
	Block
}

func (m *ConfirmReq) Read(buf *bytes.Buffer) error {
	err := m.Header.ReadHeader(buf)
	if err != nil {
		return err
	}

	if m.Header.Type != msgConfirmReq {
		return errors.New("Tried to read wrong message type")
	}
	err = m.Block.Read(m.Header.BlockType, buf)
	if err != nil {
		return err
	}

	return nil
}

func (m *ConfirmReq) Write(buf *bytes.Buffer) error {
	err := m.Header.WriteHeader(buf)
	if err != nil {
		return err
	}

	err = m.Block.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

type Publish struct {
	Header
	Block
}

func (m *Publish) Read(buf *bytes.Buffer) error {
	err := m.Header.ReadHeader(buf)
	if err != nil {
		return err
	}

	if m.Header.Type != msgPublish {
		return errors.New("Tried to read wrong message type")
	}
	err = m.Block.Read(m.Header.BlockType, buf)
	if err != nil {
		return err
	}

	return nil
}

func (m *Publish) Write(buf *bytes.Buffer) error {
	err := m.Header.WriteHeader(buf)
	if err != nil {
		return err
	}

	err = m.Block.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
