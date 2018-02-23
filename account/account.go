package account

import (
	"encoding/hex"

	"github.com/frankh/crypto/ed25519"
	"github.com/frankh/nano/address"
	"github.com/frankh/nano/blocks"
	"github.com/frankh/nano/store"
	"github.com/frankh/nano/types"
	"github.com/frankh/nano/uint128"
	"github.com/pkg/errors"
)

type Account struct {
	privateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Head       blocks.Block
	Work       *types.Work
	PoWchan    chan types.Work
}

func (a *Account) Address() types.Account {
	return address.PubKeyToAddress(a.PublicKey)
}

func New(private string) (a Account) {
	a.PublicKey, a.privateKey = address.KeypairFromPrivateKey(private)
	account := address.PubKeyToAddress(a.PublicKey)

	open := store.FetchOpen(account)
	if open != nil {
		a.Head = open
	}

	return a
}

// Returns true if the account has prepared proof of work,
func (a *Account) HasPoW() bool {
	select {
	case work := <-a.PoWchan:
		a.Work = &work
		a.PoWchan = nil
		return true
	default:
		return false
	}
}

func (a *Account) WaitPoW() {
	for !a.HasPoW() {
	}
}

func (a *Account) WaitingForPoW() bool {
	return a.PoWchan != nil
}

func (a *Account) GeneratePowSync() error {
	err := a.GeneratePoWAsync()
	if err != nil {
		return err
	}

	a.WaitPoW()
	return nil
}

// Triggers a goroutine to generate the next proof of work.
func (a *Account) GeneratePoWAsync() error {
	if a.PoWchan != nil {
		return errors.Errorf("Already generating PoW")
	}

	a.PoWchan = make(chan types.Work)

	go func(c chan types.Work, a *Account) {
		if a.Head == nil {
			c <- blocks.GenerateWorkForHash(types.BlockHash(hex.EncodeToString(a.PublicKey)))
		} else {
			c <- blocks.GenerateWork(a.Head)
		}
	}(a.PoWchan, a)

	return nil
}

func (a *Account) GetBalance() uint128.Uint128 {
	if a.Head == nil {
		return uint128.FromInts(0, 0)
	}

	return store.GetBalance(a.Head)

}

func (a *Account) Open(source types.BlockHash, representative types.Account) (*blocks.OpenBlock, error) {
	if a.Head != nil {
		return nil, errors.Errorf("Cannot open a non empty account")
	}

	if a.Work == nil {
		return nil, errors.Errorf("No PoW")
	}

	existing := store.FetchOpen(a.Address())
	if existing != nil {
		return nil, errors.Errorf("Cannot open account, open block already exists")
	}

	send_block := store.FetchBlock(source)
	if send_block == nil {
		return nil, errors.Errorf("Could not find references send")
	}

	common := blocks.CommonBlock{
		Work:      *a.Work,
		Signature: "",
	}

	block := blocks.OpenBlock{
		source,
		representative,
		a.Address(),
		common,
	}

	block.Signature = block.Hash().Sign(a.privateKey)

	if !blocks.ValidateBlockWork(&block) {
		return nil, errors.Errorf("Invalid PoW")
	}

	a.Head = &block
	return &block, nil
}

func (a *Account) Send(destination types.Account, amount uint128.Uint128) (*blocks.SendBlock, error) {
	if a.Head == nil {
		return nil, errors.Errorf("Cannot send from empty account")
	}

	if a.Work == nil {
		return nil, errors.Errorf("No PoW")
	}

	if amount.Compare(a.GetBalance()) > 0 {
		return nil, errors.Errorf("Tried to send more than balance")
	}

	common := blocks.CommonBlock{
		Work:      *a.Work,
		Signature: "",
	}

	block := blocks.SendBlock{
		a.Head.Hash(),
		destination,
		a.GetBalance().Sub(amount),
		common,
	}

	block.Signature = block.Hash().Sign(a.privateKey)

	a.Head = &block
	return &block, nil
}

func (a *Account) Receive(source types.BlockHash) (*blocks.ReceiveBlock, error) {
	if a.Head == nil {
		return nil, errors.Errorf("Cannot receive to empty account")
	}

	if a.Work == nil {
		return nil, errors.Errorf("No PoW")
	}

	send_block := store.FetchBlock(source)

	if send_block == nil {
		return nil, errors.Errorf("Source block not found")
	}

	if send_block.Type() != blocks.Send {
		return nil, errors.Errorf("Source block is not a send")
	}

	if send_block.(*blocks.SendBlock).Destination != a.Address() {
		return nil, errors.Errorf("Send is not for this account")
	}

	common := blocks.CommonBlock{
		Work:      *a.Work,
		Signature: "",
	}

	block := blocks.ReceiveBlock{
		a.Head.Hash(),
		source,
		common,
	}

	block.Signature = block.Hash().Sign(a.privateKey)

	a.Head = &block
	return &block, nil
}

func (a *Account) Change(representative types.Account) (*blocks.ChangeBlock, error) {
	if a.Head == nil {
		return nil, errors.Errorf("Cannot change on empty account")
	}

	if a.Work == nil {
		return nil, errors.Errorf("No PoW")
	}

	common := blocks.CommonBlock{
		Work:      *a.Work,
		Signature: "",
	}

	block := blocks.ChangeBlock{
		a.Head.Hash(),
		representative,
		common,
	}

	block.Signature = block.Hash().Sign(a.privateKey)

	a.Head = &block
	return &block, nil
}
