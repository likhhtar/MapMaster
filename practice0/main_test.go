package main

import (
	"testing"
)

func TestFiller(t *testing.T) {
	b := [100]byte{}
	zero := byte('0')
	one := byte('1')
	filler(b[:], zero, one)
	// Заполнить здесь ассерт, что b содержит zero и что b содержит one
	containsZero := false
	containsOne := false

	for _, v := range b {
		if v == zero {
			containsZero = true
		}
		if v == one {
			containsOne = true
		}
	}

	// Ассерты с выводом информации об ошибке
	if !containsZero {
		t.Errorf("Array b does not contain the byte zero ('%c')", zero)
	}

	if !containsOne {
		t.Errorf("Array b does not contain the byte one ('%c')", one)
	}
}
