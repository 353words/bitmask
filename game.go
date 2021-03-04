package bitmask

import (
	"fmt"
	"strings"
)

// Key is a key in the game
type Key byte

const (
	Copper  Key = 1 << iota // 1
	Jade                    // 2
	Crystal                 // 4
	maxKey
)

// String implements the fmt.Stringer interface
func (k Key) String() string {
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
	Keys Key
}

// AddKey adds a key to the player keys
func (p *Player) AddKey(key Key) {
	p.Keys |= key
}

// HasKey returns true if player has a key
func (p *Player) HasKey(key Key) bool {
	return p.Keys&key != 0
}

// RemoveKey removes key from player
func (p *Player) RemoveKey(key Key) {
	p.Keys &= ^key
}
