package bitmask

import (
	"fmt"
	"strings"
)

// KeySet is a set of keys in the game
type KeySet byte

const (
	Copper  KeySet = 1 << iota // 1
	Jade                       // 2
	Crystal                    // 4
	maxKey
)

// String implements the fmt.Stringer interface
func (k KeySet) String() string {
	if k >= maxKey {
		return fmt.Sprintf("<unknown key: %d>", k)
	}

	switch k {
	case Copper:
		return "copper"
	case Jade:
		return "jade"
	case Crystal:
		return "crystal"
	}

	// multiple keys
	var names []string
	for key := Copper; key < maxKey; key <<= 1 {
		if k&key != 0 {
			names = append(names, key.String())
		}
	}
	return strings.Join(names, "|")
}

// Player is a player in the game
type Player struct {
	Name string
	Keys KeySet
}

// AddKey adds a key to the player keys
func (p *Player) AddKey(key KeySet) {
	p.Keys |= key
}

// HasKey returns true if player has a key
func (p *Player) HasKey(key KeySet) bool {
	return p.Keys&key != 0
}

// RemoveKey removes key from player
func (p *Player) RemoveKey(key KeySet) {
	p.Keys &= ^key
}
