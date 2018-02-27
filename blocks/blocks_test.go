package blocks

import (
	"encoding/hex"
	"testing"

	"github.com/frankh/nano/address"
	"github.com/frankh/nano/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignMessage(t *testing.T) {
	testPrv := "34F0A37AAD20F4A260F0A5B3CB3D7FB50673212263E58A380BC10474BB039CE4"

	block, _ := FromJson([]byte(`{
		"type":           "open",
		"source":         "B0311EA55708D6A53C75CDBF88300259C6D018522FE3D4D0A242E431F9E8B6D0",
		"representative": "xrb_3e3j5tkog48pnny9dmfzj1r16pg8t1e76dz5tmac6iq689wyjfpiij4txtdo",
		"account":        "xrb_3e3j5tkog48pnny9dmfzj1r16pg8t1e76dz5tmac6iq689wyjfpiij4txtdo",
		"work":           "9680625b39d3363d",
		"signature":      "ECDA914373A2F0CA1296475BAEE40500A7F0A7AD72A5A80C81D7FAB7F6C802B2CC7DB50F5DD0FB25B2EF11761FA7344A158DD5A700B21BD47DE5BD0F63153A02"
	}`))

	sig := SignMessage(testPrv, block.Hash().Slice())
	assert.EqualValues(t, sig, block.GetSignature())
}

func TestGenerateWork(t *testing.T) {
	WorkThreshold = 0xfff0000000000000
	GenerateWork(LiveGenesisBlock)
}

func BenchmarkGenerateWork(b *testing.B) {
	WorkThreshold = 0xfff0000000000000
	for n := 0; n < b.N; n++ {
		GenerateWork(LiveGenesisBlock)
	}
}

func TestValidateWork(t *testing.T) {
	WorkThreshold = 0xffffffc000000000

	liveBlockHash, _ := address.AddressToPub(LiveGenesisBlock.Account)
	liveWorkBytes, _ := hex.DecodeString(string(LiveGenesisBlock.Work))
	liveBadWork, _ := hex.DecodeString("0000000000000000")

	var lbh [32]byte
	copy(lbh[:], liveBlockHash)

	require.True(t, ValidateBlockWork(LiveGenesisBlock))
	// A bit of a redundandy test to ensure ValidateBlockWork is correct
	assert.True(t, ValidateWork(lbh, utils.Reversed(liveWorkBytes)))
	assert.False(t, ValidateWork(lbh, liveBadWork))
}

func TestHashOpen(t *testing.T) {
	assert.Equal(t, LiveGenesisBlockHash, LiveGenesisBlock.Hash())
}
