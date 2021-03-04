package main

import (
	"testing"
)

var (
	mKeys = map[string]bool{
		"crystal": true,
		"copper":  true,
	}
	sKeys = []string{"crystal", "copper"}
	bKeys = Copper | Crystal
)

/*
type Key byte

const (
	Jade    Key = 1 << iota
	Copper      // 2
	Crystal     // 4
)
*/

func mHasKey(key string) bool {
	return mKeys[key]
}

func cHasKey(key string) bool {
	for _, k := range sKeys {
		if key == k {
			return true
		}
	}
	return false
}

func bHasKey(key Key) bool {
	return bKeys&key != 0
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !mHasKey("crystal") {
			b.Fatal()
		}
	}
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !cHasKey("crystal") {
			b.Fatal()
		}
	}
}

func BenchmarkBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !bHasKey(Crystal) {
			b.Fatal()
		}
	}
}
