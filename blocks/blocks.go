package blocks

import (
	"encoding/json"

	"github.com/s1na/nano/types"
	"github.com/s1na/nano/types/uint128"

	"github.com/golang/crypto/blake2b"
	"github.com/pkg/errors"
)

var LiveGenesisBlockHash, _ = types.BlockHashFromString("991CF190094C00F0B68E2E5F75F6BEE95A2E0BD93CEAA4A6734DB9F19B728948")
var LiveGenesisSourceHash, _ = types.BlockHashFromString("E89208DD038FBB269987689621D52292AE9C35941A7484756ECCED92A65093BA")

var GenesisAmount uint128.Uint128 = uint128.FromInts(0xffffffffffffffff, 0xffffffffffffffff)

const TestPrivateKey string = "34F0A37AAD20F4A260F0A5B3CB3D7FB50673212263E58A380BC10474BB039CE4"

var GenesisBlock *OpenBlock
var testGenesisBlock, _ = FromJson([]byte(`{
	"type": "open",
	"source": "B0311EA55708D6A53C75CDBF88300259C6D018522FE3D4D0A242E431F9E8B6D0",
	"representative": "xrb_3e3j5tkog48pnny9dmfzj1r16pg8t1e76dz5tmac6iq689wyjfpiij4txtdo",
	"account": "xrb_3e3j5tkog48pnny9dmfzj1r16pg8t1e76dz5tmac6iq689wyjfpiij4txtdo",
	"work": "9680625b39d3363d",
	"signature": "ECDA914373A2F0CA1296475BAEE40500A7F0A7AD72A5A80C81D7FAB7F6C802B2CC7DB50F5DD0FB25B2EF11761FA7344A158DD5A700B21BD47DE5BD0F63153A02"
}`))
var TestGenesisBlock = testGenesisBlock.(*OpenBlock)

var liveGenesisBlock, _ = FromJson([]byte(`{
	"type":           "open",
	"source":         "E89208DD038FBB269987689621D52292AE9C35941A7484756ECCED92A65093BA",
	"representative": "xrb_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3",
	"account":        "xrb_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3",
	"work":           "62f05417dd3fb691",
	"signature":      "9F0C933C8ADE004D808EA1985FA746A7E95BA2A38F867640F53EC8F180BDFE9E2C1268DEAD7C2664F356E37ABA362BC58E46DBA03E523A7B5A19E4B6EB12BB02"
}`))
var LiveGenesisBlock = liveGenesisBlock.(*OpenBlock)

type BlockType string

const (
	Open    BlockType = "open"
	Receive           = "receive"
	Send              = "send"
	Change            = "change"
	Utx               = "utx"
)

type Block interface {
	Type() BlockType
	GetSignature() types.Signature
	GetWork() types.Work
	GetRoot() types.BlockHash
	Hash() types.BlockHash
	GetPrevious() types.BlockHash
}

type CommonBlock struct {
	Work      types.Work
	Signature types.Signature
	Account   types.PubKey
}

func (b *CommonBlock) GetSignature() types.Signature {
	return b.Signature
}

func (b *CommonBlock) GetWork() types.Work {
	return b.Work
}

type RawBlock struct {
	Type           BlockType
	Source         types.BlockHash
	Representative types.PubKey
	Account        types.PubKey
	Work           types.Work
	Signature      types.Signature
	Previous       types.BlockHash
	Balance        uint128.Uint128
	Destination    types.PubKey
}

func FromJson(b []byte) (Block, error) {
	var block Block
	var raw RawBlock
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}

	common := CommonBlock{
		Work:      raw.Work,
		Signature: raw.Signature,
	}

	switch raw.Type {
	case Open:
		block = &OpenBlock{
			raw.Source,
			raw.Representative,
			raw.Account,
			common,
		}
	case Send:
		block = &SendBlock{
			raw.Previous,
			raw.Destination,
			raw.Balance,
			common,
		}
	case Receive:
		block = &ReceiveBlock{
			raw.Previous,
			raw.Source,
			common,
		}
	case Change:
		block = &ChangeBlock{
			raw.Previous,
			raw.Representative,
			common,
		}
	default:
		return nil, errors.New("unknown block type")
	}

	return block, nil
}

func (b RawBlock) Hash() [32]byte {
	switch b.Type {
	case Open:
		return HashOpen(b.Source, b.Representative, b.Account)
	case Send:
		return HashSend(b.Previous, b.Destination, b.Balance)
	case Receive:
		return HashReceive(b.Previous, b.Source)
	case Change:
		return HashChange(b.Previous, b.Representative)
	default:
		panic("Unknown block type! " + b.Type)
	}
}

func (b RawBlock) HashToString() string {
	return types.BlockHash(b.Hash()).String()
}

func SignMessage(prvStr string, message []byte) (types.Signature, error) {
	key, err := types.PrvKeyFromString(prvStr)
	if err != nil {
		return types.Signature{}, nil
	}

	_, prv, err := types.KeypairFromPrvKey(key)
	if err != nil {
		return types.Signature{}, err
	}

	return prv.Sign(message), nil
}

func HashBytes(inputs ...[]byte) types.BlockHash {
	hash, err := blake2b.New(32, nil)
	if err != nil {
		panic("Unable to create hash")
	}

	for _, b := range inputs {
		hash.Write(b)
	}

	return types.BlockHashFromSlice(hash.Sum(nil))
}

func ValidateBlockWork(b Block) bool {
	return b.GetWork().Validate(b.GetRoot())
}

func GenerateWork(b Block) types.Work {
	return types.GenerateWorkForHash(b.Hash())
}
