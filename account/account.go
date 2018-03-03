package account

import (
	"encoding/json"

	"github.com/s1na/nano/types"
	"github.com/s1na/nano/uint128"
)

type Account struct {
	PrivateKey types.PrvKey
	PublicKey  types.PubKey
	Head       types.BlockHash
	Rep        types.PubKey
	Open       types.BlockHash
	Balance    uint128.Uint128
}

func NewAccount() *Account {
	a := new(Account)

	return a
}

func (a *Account) Address() string {
	return a.PublicKey.Address()
}

func (a *Account) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return string(a.Address())
	}

	return string(b)
}

func (a *Account) Sign(data []byte) types.Signature {
	return a.PrivateKey.Sign(data)
}

/*
func (a *Account) GeneratePoW() types.Work {
	var work types.Work

	if a.Head.IsZero() {
		work = types.GenerateWorkForHash(types.BlockHashFromSlice(a.PublicKey[:]))
	} else {
		work = blocks.GenerateWork(a.Head)
	}

	return work
}*/

/*
func New(private string) (Account, error) {
	a := Account{}
	var err error
	a.PublicKey, a.PrivateKey, err = types.KeypairFromPrivateKey(private)
	if err != nil {
		return a, err
	}

	open := store.FetchOpen(types.PubKey(a.PublicKey))
	if open != nil {
		a.Head = open
	}

	return a, nil
}

func (a *Account) GetBalance() uint128.Uint128 {
	if a.Head == nil {
		return uint128.FromInts(0, 0)
	}

	return store.GetBalance(a.Head)
}

func (a *Account) Open(source types.BlockHash, representative types.PubKey) (*blocks.OpenBlock, error) {
	if a.Head != nil {
		return nil, errors.Errorf("Cannot open a non empty account")
	}

	if a.Work == nil {
		return nil, errors.Errorf("No PoW")
	}

	existing := store.FetchOpen(types.PubKey(a.PublicKey))
	if existing != nil {
		return nil, errors.Errorf("Cannot open account, open block already exists")
	}

	send_block := store.FetchBlock(source)
	if send_block == nil {
		return nil, errors.Errorf("Could not find references send")
	}

	common := blocks.CommonBlock{
		Work: *a.Work,
	}

	block := blocks.OpenBlock{
		source,
		representative,
		types.PubKey(a.PublicKey),
		common,
	}

	block.Signature = a.Sign(block.Hash().Slice())

	if !blocks.ValidateBlockWork(&block) {
		return nil, errors.Errorf("Invalid PoW")
	}

	a.Head = &block
	return &block, nil
}

func (a *Account) Send(destination types.PubKey, amount uint128.Uint128) (*blocks.SendBlock, error) {
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
		Work: *a.Work,
	}

	block := blocks.SendBlock{
		a.Head.Hash(),
		destination,
		a.GetBalance().Sub(amount),
		common,
	}

	block.Signature = block.Hash().Sign(a.PrivateKey)

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

	sendBlock := store.FetchBlock(source)

	if sendBlock == nil {
		return nil, errors.Errorf("Source block not found")
	}

	if sendBlock.Type() != blocks.Send {
		return nil, errors.Errorf("Source block is not a send")
	}

	if !sendBlock.(*blocks.SendBlock).Destination.Equal(types.PubKeyFromSlice(a.PublicKey)) {
		return nil, errors.Errorf("Send is not for this account")
	}

	common := blocks.CommonBlock{
		Work: *a.Work,
	}

	block := blocks.ReceiveBlock{
		a.Head.Hash(),
		source,
		common,
	}

	block.Signature = block.Hash().Sign(a.PrivateKey)

	a.Head = &block
	return &block, nil
}

func (a *Account) Change(representative types.PubKey) (*blocks.ChangeBlock, error) {
	if a.Head == nil {
		return nil, errors.Errorf("Cannot change on empty account")
	}

	if a.Work == nil {
		return nil, errors.Errorf("No PoW")
	}

	common := blocks.CommonBlock{
		Work: *a.Work,
	}

	block := blocks.ChangeBlock{
		a.Head.Hash(),
		representative,
		common,
	}

	block.Signature = block.Hash().Sign(a.PrivateKey)

	a.Head = &block
	return &block, nil
}*/
