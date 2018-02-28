package types

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/frankh/crypto/ed25519"
	"github.com/frankh/nano/utils"

	"github.com/golang/crypto/blake2b"
	"github.com/pkg/errors"
)

const EncodeXrb = "13456789abcdefghijkmnopqrstuwxyz"

var XrbEncoding = base32.NewEncoding(EncodeXrb)

func ValidateAddress(addr string) bool {
	_, err := AccPubFromString(addr)
	return err == nil
}

func KeypairFromPrivateKey(prvStr string) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	prvBytes, _ := hex.DecodeString(prvStr)
	pub, prv, err := ed25519.GenerateKey(bytes.NewReader(prvBytes))

	return pub, prv, err
}

func KeypairFromSeed(seed string, index uint32) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	// This seems to be the standard way of producing wallets.

	// We hash together the seed with an address index and use
	// that as the private key. Whenever you "add" an address
	// to your wallet the wallet software increases the index
	// and generates a new address.
	hash, err := blake2b.New(32, nil)
	if err != nil {
		return nil, nil, err
	}

	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		return nil, nil, err
	}

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, index)

	hash.Write(seedBytes)
	hash.Write(bs)

	in := hash.Sum(nil)
	pub, prv, err := ed25519.GenerateKey(bytes.NewReader(in))
	if err != nil {
		return nil, nil, err
	}

	return pub, prv, err
}

func GenerateKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}

type AccPub ed25519.PublicKey

func AccPubFromString(addr string) (AccPub, error) {
	var a AccPub
	// A valid xrb address is 64 bytes long
	// First 4 are simply a hard-coded string xrb_ for ease of use
	// The following 52 characters form the address, and the final
	// 8 are a checksum.
	// They are base 32 encoded with a custom encoding.
	if len(addr) != 64 || addr[:4] != "xrb_" {
		return nil, errors.New("Invalid address format")
	}

	// The xrb address string is 260bits which doesn't fall on a
	// byte boundary. pad with zeros to 280bits.
	// (zeros are encoded as 1 in xrb's 32bit alphabet)
	key := "1111" + addr[4:56]
	checksum := addr[56:]

	b, err := XrbEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	// Strip off upper 24 bits (3 bytes). 20 padding was added by us,
	// 4 is unused as account is 256 bits.
	a = AccPub(b[3:])

	// Xrb checksum is calculated by hashing the key and reversing the bytes
	valid := XrbEncoding.EncodeToString(a.Checksum()) == checksum
	if !valid {
		return nil, errors.New("Invalid address checksum")
	}

	return a, nil
}

func AccPubFromSlice(data []byte) AccPub {
	var a AccPub
	copy(a[:], data)
	return a
}

func (a AccPub) Checksum() []byte {
	hash, err := blake2b.New(5, nil)
	if err != nil {
		panic("Unable to create hash")
	}

	hash.Write(a)
	return utils.Reversed(hash.Sum(nil))
}

func (a AccPub) Equal(o AccPub) bool {
	return reflect.DeepEqual(a, o)
}

func (a AccPub) String() string {
	// Pubkey is 256bits, base32 must be multiple of 5 bits
	// to encode properly.
	// Pad the start with 0's and strip them off after base32 encoding
	padded := append([]byte{0, 0, 0}, a...)
	address := XrbEncoding.EncodeToString(padded)[4:]
	checksum := XrbEncoding.EncodeToString(a.Checksum())

	return "xrb_" + address + checksum
}

func (a *AccPub) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	r, err := AccPubFromString(s)
	if err != nil {
		return err
	}

	*a = r

	return nil
}

func (a AccPub) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}
