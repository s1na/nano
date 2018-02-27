package network

import (
	"github.com/pkg/errors"
)

var (
	MagicNumber = [2]byte{'R', 'C'}
)

const (
	HeaderSize = 8
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
	invalidBlock byte = iota
	notABlock
	sendBlock
	receiveBlock
	openBlock
	changeBlock
)

type MessagePart interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type Message struct {
	Header *Header
	Body   MessagePart
}

func NewMessage(t byte, b MessagePart) *Message {
	m := new(Message)

	m.Header = NewHeader(t)
	m.Body = b

	return m
}

func (m *Message) Marshal() ([]byte, error) {
	var data []byte

	data, err := m.Header.Marshal()
	if err != nil {
		return nil, err
	}

	bb, err := m.Body.Marshal()
	if err != nil {
		return nil, err
	}
	data = append(data, bb...)

	return data, nil
}

func (m *Message) Unmarshal(data []byte) error {
	hb, bb := data[:HeaderSize], data[HeaderSize:]
	if len(hb) != HeaderSize || len(bb) == 0 {
		return errors.New("invalid message parts")
	}

	h := new(Header)
	if err := h.Unmarshal(hb); err != nil {
		return err
	}
	m.Header = h

	if m.Header.MagicNumber != MagicNumber {
		return errors.New("Header has invalid magic number")
	}

	if err := m.UnmarshalBody(bb); err != nil {
		return err
	}

	return nil
}

func (m *Message) UnmarshalBody(data []byte) error {
	switch m.Header.Type {
	case msgKeepalive:
		m.Body = new(KeepAlive)
	case msgPublish:
		m.Body = &Publish{Block: Block{Type: m.Header.BlockType}}
	case msgConfirmReq:
		m.Body = &ConfirmReq{Block: Block{Type: m.Header.BlockType}}
	case msgConfirmAck:
		m.Body = &ConfirmAck{Vote: Vote{Block: Block{Type: m.Header.BlockType}}}
	default:
		return errors.New("message type undefined")
	}

	if err := m.Body.Unmarshal(data); err != nil {
		return errors.Wrap(err, "failed to unmarshal message body")
	}

	return nil
}

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

func (h *Header) Marshal() ([]byte, error) {
	data := []byte{
		h.MagicNumber[0],
		h.MagicNumber[1],
		h.VersionMax,
		h.VersionUsing,
		h.VersionMin,
		h.Type,
		h.Extensions,
		h.BlockType,
	}

	return data, nil
}

func (h *Header) Unmarshal(data []byte) error {
	if len(data) != 8 {
		return errors.New("header len is invalid")
	}

	h.MagicNumber[0] = data[0]
	h.MagicNumber[1] = data[1]
	h.VersionMax = data[2]
	h.VersionUsing = data[3]
	h.VersionMin = data[4]
	h.Type = data[5]
	h.Extensions = data[6]
	h.BlockType = data[7]

	return nil
}

type KeepAlive struct {
	Peers [8]Peer
}

func NewKeepAlive(peers []Peer) *KeepAlive {
	m := new(KeepAlive)

	for i, p := range peers {
		if i == len(m.Peers) {
			break
		}

		m.Peers[i] = p
	}

	return m
}

func (m *KeepAlive) Unmarshal(data []byte) error {
	if len(data) != len(m.Peers)*(16+2) {
		return errors.New("keepalive packet has invalid length")
	}

	for i := 0; i < len(m.Peers); i++ {
		p := NewPeer()
		if err := p.Unmarshal(data[i*18 : (i+1)*18]); err != nil {
			return err
		}

		m.Peers[i] = *p
	}

	return nil
}

func (m *KeepAlive) Marshal() ([]byte, error) {
	data := make([]byte, 0, 8*(16+2))

	for _, peer := range m.Peers {
		pb, err := peer.Marshal()
		if err != nil {
			return nil, err
		}

		data = append(data, pb...)
	}

	return data, nil
}

type Publish struct {
	Block
}

type ConfirmReq struct {
	Block
}

type ConfirmAck struct {
	Vote
}
