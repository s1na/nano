package store

import (
	"github.com/dgraph-io/badger"
)

var (
	store *Store
)

type Store struct {
	db      *badger.DB
	dataDir string
}

func NewStore(dataDir string) *Store {
	s := new(Store)

	s.dataDir = dataDir

	return s
}

func (s *Store) Start() error {
	opts := badger.DefaultOptions
	opts.Dir = s.dataDir
	opts.ValueDir = s.dataDir
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Stop() {
	s.db.Close()
}

func (s *Store) Set(k []byte, v []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(k, v); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) Get(k []byte) ([]byte, error) {
	var v []byte
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(k)
	if err != nil {
		return nil, err
	}

	v, err = item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Store) GetKeys() [][]byte {
	keys := make([][]byte, 0, 2)
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	it := txn.NewIterator(opts)
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		keys = append(keys, item.Key())
	}

	return keys
}

func (s *Store) GetPrefixKeys(prefix []byte) [][]byte {
	keys := make([][]byte, 0, 2)
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		keys = append(keys, item.Key())
	}

	return keys
}

func (s *Store) GetPrefixValues(prefix []byte) (map[string][]byte, error) {
	res := make(map[string][]byte)
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		v, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		res[string(k)] = v
	}

	return res, nil
}

// Blocks that we cannot store due to not having their parent
// block stored
/*var unconnectedBlockPool map[types.BlockHash]blocks.Block

func Init(config Config) {
	var err error
	unconnectedBlockPool = make(map[types.BlockHash]blocks.Block)

	if globalConn != nil {
		globalConn.Close()
		globalConn = nil
	}
	Conf = &config
	conn := getConn()
	defer releaseConn(conn)

	_, err = conn.Get(blocks.LiveGenesisBlockHash[:])

	if err != nil {
		uncheckedStoreBlock(conn, config.GenesisBlock)
	}
}*/

/*
func GetBalance(block blocks.Block) uint128.Uint128 {
	conn := getConn()
	defer releaseConn(conn)
	return getBalance(conn, block)
}

func getSendAmount(conn *badger.Txn, block *blocks.SendBlock) uint128.Uint128 {
	prev := fetchBlock(conn, block.PreviousHash)

	return getBalance(conn, prev).Sub(getBalance(conn, block))
}

func getBalance(conn *badger.Txn, block blocks.Block) uint128.Uint128 {
	switch block.Type() {
	case blocks.Open:
		b := block.(*blocks.OpenBlock)
		if b.SourceHash == Conf.GenesisBlock.SourceHash {
			return blocks.GenesisAmount
		}
		source := fetchBlock(conn, b.SourceHash).(*blocks.SendBlock)
		return getSendAmount(conn, source)

	case blocks.Send:
		b := block.(*blocks.SendBlock)
		return b.Balance

	case blocks.Receive:
		b := block.(*blocks.ReceiveBlock)
		prev := fetchBlock(conn, b.PreviousHash)
		source := fetchBlock(conn, b.SourceHash).(*blocks.SendBlock)
		received := getSendAmount(conn, source)
		return getBalance(conn, prev).Add(received)

	case blocks.Change:
		b := block.(*blocks.ChangeBlock)
		return getBalance(conn, fetchBlock(conn, b.PreviousHash))

	default:
		panic("Unknown block type")
	}
}
*/
/*
func StoreBlock(block blocks.Block) error {
	conn := getConn()
	defer releaseConn(conn)
	return storeBlock(conn, block)
}

func storeBlock(conn *badger.Txn, block blocks.Block) error {
	if !blocks.ValidateBlockWork(block) {
		return errors.New("Invalid work for block")
	}

	if block.Type() != blocks.Open && block.Type() != blocks.Change && block.Type() != blocks.Send && block.Type() != blocks.Receive {
		return errors.New("Unknown block type")
	}

	if fetchBlock(conn, block.PreviousBlockHash()) == nil {
		if unconnectedBlockPool[block.PreviousBlockHash()] == nil {
			unconnectedBlockPool[block.PreviousBlockHash()] = block
			log.Printf("Added block to unconnected pool, now %d", len(unconnectedBlockPool))
		}
		return errors.New("Cannot find parent block")
	}

	uncheckedStoreBlock(conn, block)
	dependentBlock := unconnectedBlockPool[block.Hash()]

	if dependentBlock != nil {
		// We have an unconnected block dependent on this: Store it now that
		// it's connected
		delete(unconnectedBlockPool, block.Hash())
		storeBlock(conn, dependentBlock)
	}

	return nil
}

func uncheckedStoreBlock(conn *badger.Txn, block blocks.Block) {
	var buf bytes.Buffer
	var meta byte
	enc := gob.NewEncoder(&buf)
	switch block.Type() {
	case blocks.Open:
		b := block.(*blocks.OpenBlock)
		meta = MetaOpen
		err := enc.Encode(b)
		if err != nil {
			panic(err)
		}
		// Open blocks need to be stored twice, once keyed on account,
		// once keyed on hash.
		err = conn.SetWithMeta(b.RootHash().Slice(), buf.Bytes(), meta)
		if err != nil {
			panic(err)
		}
	case blocks.Send:
		b := block.(*blocks.SendBlock)
		meta = MetaSend
		err := enc.Encode(b)
		if err != nil {
			panic(err)
		}
	case blocks.Receive:
		b := block.(*blocks.ReceiveBlock)
		meta = MetaReceive
		err := enc.Encode(b)
		if err != nil {
			panic(err)
		}
	case blocks.Change:
		b := block.(*blocks.ChangeBlock)
		meta = MetaChange
		err := enc.Encode(b)
		if err != nil {
			panic(err)
		}
	default:
		panic("Unknown block type")
	}

	err := conn.SetWithMeta(block.Hash().Slice(), buf.Bytes(), meta)
	if err != nil {
		panic("Failed to store block")
	}
}*/
