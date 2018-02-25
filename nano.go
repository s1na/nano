package main

import (
	"math/rand"
	"time"

	"github.com/frankh/nano/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
