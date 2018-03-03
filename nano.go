package main

import (
	"math/rand"
	"time"

	"github.com/s1na/nano/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
