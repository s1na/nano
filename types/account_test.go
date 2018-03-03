package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validAddresses = []string{
	"xrb_38nm8t5rimw6h6j7wyokbs8jiygzs7baoha4pqzhfw1k79npyr1km8w6y7r8",
	"xrb_1awsn43we17c1oshdru4azeqjz9wii41dy8npubm4rg11so7dx3jtqgoeahy",
	"xrb_3arg3asgtigae3xckabaaewkx3bzsh7nwz7jkmjos79ihyaxwphhm6qgjps4",
	"xrb_3pczxuorp48td8645bs3m6c3xotxd3idskrenmi65rbrga5zmkemzhwkaznh",
	"xrb_3hd4ezdgsp15iemx7h81in7xz5tpxi43b6b41zn3qmwiuypankocw3awes5k",
	"xrb_1anrzcuwe64rwxzcco8dkhpyxpi8kd7zsjc1oeimpc3ppca4mrjtwnqposrs",
}

var invalidAddresses = []string{
	"xrb_38nm8t5rimw6h6j7wyokbs8jiygzs7baoha4pqzhfw1k79npyr1km8w6y7r7",
	"xrc_38nm8t5rimw6h6j7wyokbs8jiygzs7baoha4pqzhfw1k79npyr1km8w6y7r8",
	"xrb38nm8t5rimw6h6j7wyokbs8jiygzs7baoha4pqzhfw1k79npyr1km8w6y7r8",
	"xrb8nm8t5rimw6h6j7wyokbs8jiygzs7baoha4pqzhfw1k79npyr1km8w6y7r8",
	"xrb_8nm8t5rimw6h6j7wyokbs8jiygzs7baoha4pqzhfw1k79npyr1km8w6y7r8",
}

func TestPubKeyFromAddress(t *testing.T) {
	a, err := PubKeyFromAddress("xrb_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3")
	require.Nil(t, err)
	assert.Equal(t, "e89208dd038fbb269987689621d52292ae9c35941a7484756ecced92a65093ba", hex.EncodeToString(a))
}

func TestValidateAddress(t *testing.T) {
	for _, addr := range validAddresses {
		assert.True(t, ValidateAddress(addr))
	}

	for _, addr := range invalidAddresses {
		assert.False(t, ValidateAddress(addr))
	}
}

func TestKeypairFromSeed(t *testing.T) {
	seed := "1234567890123456789012345678901234567890123456789012345678901234"
	seedBytes, err := hex.DecodeString(seed)
	require.Nil(t, err)

	// Generated from the official RaiBlocks wallet using above seed.
	expected := map[uint32]string{
		0: "xrb_3iwi45me3cgo9aza9wx5f7rder37hw11xtc1ek8psqxw5oxb8cujjad6qp9y",
		1: "xrb_3a9d1h6wt3zp8cqd6dhhgoyizmk1ciemqkrw97ysrphn7anm6xko1wxakaa1",
		2: "xrb_1dz36wby1azyjgh7t9nopjm3k5rduhmntercoz545my9s8nm7gcuthuq9fmq",
		3: "xrb_1fb7kaqaue49kf9w4mb9w3scuxipbdm3ez6ibnri4w8qexzg5f4r7on1dmxb",
		4: "xrb_3h9a64yqueuij1j9odt119r3ymm8n83wyyz7o9u7ram1tgfhsh1zqwjtzid9",
	}

	for k, v := range expected {
		pub, _, err := KeypairFromSeed(PrvKey(seedBytes), k)
		require.Nil(t, err)
		assert.Equal(t, v, PubKey(pub).Address())
	}
}

func TestKeypairFromPrvKey(t *testing.T) {
	keyStr := "9C5CC0CA478CE841E33806D05100644C592DAA825E6AC13A125D0731AD0424C51480F0DD781A863189DA84970CB717B956984EFA931890677F10C04337ADE4FD"
	key, err := PrvKeyFromString(keyStr)
	require.Nil(t, err)

	_, prv, err := KeypairFromPrvKey(key)
	require.Nil(t, err)

	assert.Equal(t, key, prv)
}

func BenchmarkGenerateAddress(b *testing.B) {
	for n := 0; n < b.N; n++ {
		pub, _, _ := GenerateKey(nil)
		PubKey(pub).Address()
	}
}
