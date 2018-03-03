package ledger

import (
	"os"
	"testing"

	"github.com/s1na/nano/account"
	"github.com/s1na/nano/blocks"
	"github.com/s1na/nano/store"
	"github.com/s1na/nano/types"
	"github.com/s1na/nano/types/uint128"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LedgerTestSuite struct {
	suite.Suite
	st *store.Store
	bs *blocks.BlockStore
	as *account.AccountStore
}

func (s *LedgerTestSuite) SetupTest() {
	blocks.GenesisBlock = blocks.TestGenesisBlock
	types.WorkThreshold = uint64(0xff00000000000000)

	s.st = store.NewStore("testdata")
	s.st.Start()
	s.bs = blocks.NewBlockStore(s.st)
	s.as = account.NewAccountStore(s.st)
}

func (s *LedgerTestSuite) TearDownTest() {
	s.st.Stop()
	os.RemoveAll("testdata")
}

func (s *LedgerTestSuite) TestInit() {
	l := NewLedger(s.st)
	err := l.Init()
	require.Nil(s.T(), err)

	b, err := s.bs.GetBlock(blocks.GenesisBlock.Hash())
	require.Nil(s.T(), err)

	ob := b.(*blocks.OpenBlock)
	s.Equal(blocks.GenesisBlock, ob)

	acc, err := s.as.GetAccount(ob.Account)
	require.Nil(s.T(), err)
	s.EqualValues(blocks.GenesisBlock.Account, acc.PublicKey)
	s.Equal(blocks.GenesisAmount, acc.Balance)
}

func (s *LedgerTestSuite) TestAddSend() {
	l := NewLedger(s.st)
	err := l.Init()
	require.Nil(s.T(), err)

	dest, _, err := types.GenerateKey(nil)
	require.Nil(s.T(), err)

	b := &blocks.SendBlock{
		Previous:    blocks.TestGenesisBlock.Hash(),
		Destination: types.PubKeyFromSlice(dest),
		Balance:     blocks.GenesisAmount.Sub(uint128.FromInts(0, 1000)),
	}
	b.Work = types.GenerateWorkForHash(b.GetRoot())
	err = l.AddSend(b)
	require.Nil(s.T(), err)

	r, err := s.bs.GetBlock(b.Hash())
	require.Nil(s.T(), err)

	sb, ok := r.(*blocks.SendBlock)
	require.True(s.T(), ok)
	s.Equal(b, sb)
}

func TestLedgerTestSuite(t *testing.T) {
	suite.Run(t, new(LedgerTestSuite))
}
