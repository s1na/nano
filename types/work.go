package types

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/frankh/nano/utils"

	"github.com/golang/crypto/blake2b"
)

var WorkThreshold = uint64(0xffffffc000000000)

func GenerateWorkForHash(b BlockHash) Work {
	work := Work{0, 0, 0, 0, 0, 0, 0, 0}
	for {
		if work.Validate(b) {
			return work
		}

		incrementWork(work[:])
	}
}

func incrementWork(work []byte) {
	for i := 0; i < len(work)-1; i++ {
		if work[i] < 255 {
			work[i]++
			return
		} else {
			work[i]++
			incrementWork(work[i+1:])
			return
		}
	}
}

type Work [8]byte

func WorkFromString(s string) (Work, error) {
	var w Work
	b, err := hex.DecodeString(s)
	if err != nil {
		return w, err
	}

	copy(w[:], utils.Reversed(b))

	return w, nil
}

func WorkFromSlice(b []byte) Work {
	var w Work
	copy(w[:], b)

	return w
}

// ValidateWork takes the "work" value (little endian from hex)
// and block hash and verifies that the work passes the difficulty.
// To verify this, we create a new 8 byte hash of the
// work and the block hash and convert this to a uint64
// which must be higher (or equal) than the difficulty
// (0xffffffc000000000) to be valid.
func (w Work) Validate(blockHash BlockHash) bool {
	hash, err := blake2b.New(8, nil)
	if err != nil {
		panic("Unable to create hash")
	}

	hash.Write(w[:])
	hash.Write(blockHash[:])

	workValue := hash.Sum(nil)
	workValueInt := binary.LittleEndian.Uint64(workValue)

	return workValueInt >= WorkThreshold
}

func (w Work) Slice() []byte {
	return w[:]
}

func (w Work) String() string {
	return strings.ToUpper(hex.EncodeToString(utils.Reversed(w[:])))
}

func (w *Work) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	b, err := WorkFromString(s)
	if err != nil {
		return err
	}

	*w = b

	return nil
}

func (w Work) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}
