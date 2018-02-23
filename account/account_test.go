package account

import (
	"encoding/hex"
	"testing"

	"github.com/frankh/nano/address"
	"github.com/frankh/nano/blocks"
	"github.com/frankh/nano/store"
	"github.com/frankh/nano/uint128"
)

func TestNew(t *testing.T) {
	store.Init(store.TestConfig)

	a := New(blocks.TestPrivateKey)
	if a.GetBalance() != blocks.GenesisAmount {
		t.Errorf("Genesis block doesn't have correct balance")
	}
}

func TestPoW(t *testing.T) {
	blocks.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	a := New(blocks.TestPrivateKey)

	if a.GeneratePoWAsync() != nil || !a.WaitingForPoW() {
		t.Errorf("Failed to start PoW generation")
	}

	if a.GeneratePoWAsync() == nil {
		t.Errorf("Started PoW while already in progress")
	}

	_, err := a.Send(blocks.TestGenesisBlock.Account, uint128.FromInts(0, 1))

	if err == nil {
		t.Errorf("Created send block without PoW")
	}

	a.WaitPoW()

	send, _ := a.Send(blocks.TestGenesisBlock.Account, uint128.FromInts(0, 1))

	if !blocks.ValidateBlockWork(send) {
		t.Errorf("Invalid work")
	}

}

func TestSend(t *testing.T) {
	blocks.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	a := New(blocks.TestPrivateKey)

	a.GeneratePowSync()
	amount := uint128.FromInts(1, 1)

	send, _ := a.Send(blocks.TestGenesisBlock.Account, amount)

	if a.GetBalance() != blocks.GenesisAmount.Sub(amount) {
		t.Errorf("Balance unchanged after send")
	}

	_, err := a.Send(blocks.TestGenesisBlock.Account, blocks.GenesisAmount)
	if err == nil {
		t.Errorf("Sent more than account balance")
	}

	a.GeneratePowSync()
	store.StoreBlock(send)
	receive, _ := a.Receive(send.Hash())
	store.StoreBlock(receive)

	if a.GetBalance() != blocks.GenesisAmount {
		t.Errorf("Balance not updated after receive, %x != %x", a.GetBalance().GetBytes(), blocks.GenesisAmount.GetBytes())
	}

}

func TestOpen(t *testing.T) {
	blocks.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	amount := uint128.FromInts(1, 1)

	sendW := New(blocks.TestPrivateKey)
	sendW.GeneratePowSync()

	_, priv := address.GenerateKey()
	openW := New(hex.EncodeToString(priv))
	send, _ := sendW.Send(openW.Address(), amount)
	openW.GeneratePowSync()

	_, err := openW.Open(send.Hash(), openW.Address())
	if err == nil {
		t.Errorf("Expected error for referencing unstored send")
	}

	if openW.GetBalance() != uint128.FromInts(0, 0) {
		t.Errorf("Open should start at zero balance")
	}

	store.StoreBlock(send)
	_, err = openW.Open(send.Hash(), openW.Address())
	if err != nil {
		t.Errorf("Open block failed: %s", err)
	}

	if openW.GetBalance() != amount {
		t.Errorf("Open balance didn't equal send amount")
	}

	_, err = openW.Open(send.Hash(), openW.Address())
	if err == nil {
		t.Errorf("Expected error for creating duplicate open block")
	}

}
