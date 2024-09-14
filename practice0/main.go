package main

import (
	"math/rand"
	"time"
)

func main() {
}

func filler(b []byte, ifzero byte, ifnot byte) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(b); i++ {
		if rand.Intn(2) == 0 {
			b[i] = ifzero
		} else {
			b[i] = ifnot
		}
	}
}
