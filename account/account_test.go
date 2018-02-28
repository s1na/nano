package account

import (
/*	"testing"

	"github.com/frankh/nano/blocks"
	"github.com/frankh/nano/store"
	"github.com/frankh/nano/types"
	"github.com/frankh/nano/uint128"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"*/
)

/*
func TestPoW(t *testing.T) {
	types.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	a, err := New(blocks.TestPrivateKey)

	require.Nil(t, err)
	if a.GeneratePoWAsync() != nil || !a.WaitingForPoW() {
		t.Errorf("Failed to start PoW generation")
	}

	assert.NotNil(t, a.GeneratePoWAsync())

	_, err = a.Send(blocks.TestGenesisBlock.Account, uint128.FromInts(0, 1))
	require.NotNil(t, err)

	a.WaitPoW()

	send, _ := a.Send(blocks.TestGenesisBlock.Account, uint128.FromInts(0, 1))
	assert.True(t, blocks.ValidateBlockWork(send))
}
*/
/*
func TestSend(t *testing.T) {
	blocks.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	a, err := New(blocks.TestPrivateKey)
	require.Nil(t, err)

	a.GeneratePowSync()
	amount := uint128.FromInts(1, 1)

	send, err := a.Send(blocks.TestGenesisBlock.Account, amount)
	require.Nil(t, err)
	assert.Equal(t, blocks.GenesisAmount.Sub(amount), a.GetBalance())

	_, err = a.Send(blocks.TestGenesisBlock.Account, blocks.GenesisAmount)
	assert.NotNil(t, err)

	a.GeneratePowSync()
	store.StoreBlock(send)
	receive, err := a.Receive(send.Hash())
	require.Nil(t, err)
	store.StoreBlock(receive)

	assert.Equal(t, blocks.GenesisAmount, a.GetBalance())
}

func TestOpen(t *testing.T) {
	blocks.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	amount := uint128.FromInts(1, 1)

	sendW, err := New(blocks.TestPrivateKey)
	require.Nil(t, err)
	sendW.GeneratePowSync()

	_, priv, err := types.GenerateKey()
	require.Nil(t, err)

	openW, err := New(hex.EncodeToString(priv))
	require.Nil(t, err)

	send, _ := sendW.Send(types.AccPubFromSlice(openW.PublicKey), amount)
	openW.GeneratePowSync()

	_, err = openW.Open(send.Hash(), types.AccPubFromSlice(openW.PublicKey))
	if err == nil {
		t.Errorf("Expected error for referencing unstored send")
	}

	if openW.GetBalance() != uint128.FromInts(0, 0) {
		t.Errorf("Open should start at zero balance")
	}

	store.StoreBlock(send)
	_, err = openW.Open(send.Hash(), types.AccPubFromSlice(openW.PublicKey))
	if err != nil {
		t.Errorf("Open block failed: %s", err)
	}

	if openW.GetBalance() != amount {
		t.Errorf("Open balance didn't equal send amount")
	}

	_, err = openW.Open(send.Hash(), types.AccPubFromSlice(openW.PublicKey))
	if err == nil {
		t.Errorf("Expected error for creating duplicate open block")
	}

}*/
