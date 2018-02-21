package node

import (
	"bytes"
	"testing"

	"github.com/frankh/nano/store"
)

func TestHandleMessage(t *testing.T) {
	store.Init(store.TestConfig)
	NewNetwork().handleMessage("::1", bytes.NewBuffer(publishTest))
}
