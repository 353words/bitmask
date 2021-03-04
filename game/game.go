package main

import (
	"fmt"
	"strings"
)

type Player struct {
	Name string
	Keys Key
}

func (p *Player) AddKey(key Key) {
	p.Keys |= key
}

func (p *Player) HasKey(key Key) bool {
	return p.Keys&key != 0
}

func (p *Player) RemoveKey(key Key) {
	p.Keys &= ^key
}

// Key is a key in the game
type Key byte

const (
	Copper  Key = 1 << iota // 1
	Jade                    // 2
	Crystal                 // 4
	MaxKey
)

// golang.org/x/tools/cmd/stringer
// implement fmt.Stringer interface
func (k Key) String() string {
	if k >= MaxKey {
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
	for key := Copper; key < MaxKey; key <<= 1 {
		if k&key != 0 {
			names = append(names, key.String())
		}
	}
	return strings.Join(names, "|")
}

func main() {
	p := Player{"Parzival", 0}
	fmt.Printf("%+v\n", p)
	p.AddKey(Copper)
	fmt.Printf("%+v\n", p)
	p.AddKey(Jade)
	fmt.Printf("%+v\n", p)
	p.RemoveKey(Copper)
	fmt.Printf("%+v\n", p)
	fmt.Println(p.HasKey(Copper))
	fmt.Println(p.HasKey(Jade))
}
