package types

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/frankh/crypto/ed25519"

	"github.com/golang/crypto/blake2b"
	"github.com/pkg/errors"
)

const EncodeXrb = "13456789abcdefghijkmnopqrstuwxyz"

var XrbEncoding = base32.NewEncoding(EncodeXrb)

func ValidateAddress(addr string) bool {
	_, err := PubKeyFromAddress(addr)
	return err == nil
}

func GenerateKey(rd io.Reader) (PubKey, PrvKey, error) {
	pub, prv, err := ed25519.GenerateKey(rd)
	return PubKey(pub), PrvKey(prv), err
}

func KeypairFromPrvKey(key PrvKey) (PubKey, PrvKey, error) {
	return GenerateKey(bytes.NewReader(key))
}

func KeypairFromSeed(seed PrvKey, index uint32) (PubKey, PrvKey, error) {
	// This seems to be the standard way of producing wallets.

	// We hash together the seed with an address index and use
	// that as the private key. Whenever you "add" an address
	// to your wallet the wallet software increases the index
	// and generates a new address.
	hash, err := blake2b.New(32, nil)
	if err != nil {
		return nil, nil, err
	}

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, index)

	hash.Write(seed)
	hash.Write(bs)

	in := hash.Sum(nil)
	return GenerateKey(bytes.NewReader(in))
}

type PubKey ed25519.PublicKey

func PubKeyFromAddress(addr string) (PubKey, error) {
	var a PubKey
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
	a = PubKey(b[3:])

	// Xrb checksum is calculated by hashing the key and reversing the bytes
	valid := XrbEncoding.EncodeToString(a.Checksum()) == checksum
	if !valid {
		return nil, errors.New("Invalid address checksum")
	}

	return a, nil
}

func PubKeyFromHex(s string) (PubKey, error) {
	b, err := hex.DecodeString(s)
	return PubKey(b), err
}

func PubKeyFromSlice(data []byte) PubKey {
	var a PubKey
	copy(a[:], data)
	return a
}

func (a PubKey) Checksum() []byte {
	hash, err := blake2b.New(5, nil)
	if err != nil {
		panic("Unable to create hash")
	}

	hash.Write(a)
	return Reversed(hash.Sum(nil))
}

func (a PubKey) Equal(o PubKey) bool {
	return reflect.DeepEqual(a, o)
}

func (a PubKey) Address() string {
	// Pubkey is 256bits, base32 must be multiple of 5 bits
	// to encode properly.
	// Pad the start with 0's and strip them off after base32 encoding
	padded := append([]byte{0, 0, 0}, a...)
	address := XrbEncoding.EncodeToString(padded)[4:]
	checksum := XrbEncoding.EncodeToString(a.Checksum())

	return "xrb_" + address + checksum
}

func (a PubKey) Hex() string {
	return strings.ToUpper(hex.EncodeToString(a))
}

func (a *PubKey) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// TODO: determine representation (hex, addr?)
	r, err := PubKeyFromAddress(s)
	if err != nil {
		return err
	}

	*a = r

	return nil
}

func (a PubKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Address())
}

type PrvKey ed25519.PrivateKey

func PrvKeyFromString(data string) (PrvKey, error) {
	b, err := hex.DecodeString(data)
	return PrvKey(b), err
}

func PrvKeyFromSlice(data []byte) PrvKey {
	return PrvKey(data)
}

func (k PrvKey) Sign(data []byte) Signature {
	return SignatureFromSlice(ed25519.Sign(ed25519.PrivateKey(k), data))
}

func (k PrvKey) Equal(o PrvKey) bool {
	return reflect.DeepEqual(k, o)
}

func (k PrvKey) String() string {
	return strings.ToUpper(hex.EncodeToString(k))
}

func (k *PrvKey) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	r, err := PrvKeyFromString(s)
	if err != nil {
		return err
	}

	*k = r

	return nil
}

func (k PrvKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}
