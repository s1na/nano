package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func NewPeer() *Peer {
	p := new(Peer)
	return p
}

func (p *Peer) Addr() *net.UDPAddr {
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.IP.String(), p.Port))
	return addr
}

func (p *Peer) String() string {
	return fmt.Sprintf("%s:%d", p.IP.String(), p.Port)
}

func (p *Peer) Marshal() ([]byte, error) {
	data := make([]byte, net.IPv6len+2)

	copy(data[:net.IPv6len], p.IP[:])
	binary.LittleEndian.PutUint16(data[net.IPv6len:], p.Port)

	return data, nil
}

func (p *Peer) Unmarshal(data []byte) error {
	if len(data) != net.IPv6len+2 {
		return errors.New("peer to be unmarshalled has invalid length")
	}

	p.IP = net.IP(data[:net.IPv6len])
	p.Port = binary.LittleEndian.Uint16(data[net.IPv6len:])

	return nil
}
