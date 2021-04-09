package bitmask

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

func bHasKey(key KeySet) bool {
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

func BenchmarkMemory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sl := make([]string, 2)
		sl[0] = "copper"
		sl[1] = "jade"
		if len(sl) != 2 {
			b.Fatal(sl)
		}

	}
}
