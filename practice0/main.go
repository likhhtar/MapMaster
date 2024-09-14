package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	b := make([]byte, 100)
	go func() {
		for {
			filler(b[:50], '0', '1')
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			filler(b[50:], 'X', 'Y')
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			fmt.Println(string(b))
			time.Sleep(time.Second)
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}
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
