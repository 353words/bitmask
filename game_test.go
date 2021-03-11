package bitmask

import "testing"

func TestKeys(t *testing.T) {
	p := Player{"Parzival", 0}
	if p.Keys.String() != "" {
		t.Fatalf("empty keys: %q", p.Keys)
	}

	p.AddKey(Copper)
	if p.Keys.String() != Copper.String() {
		t.Fatalf("+copper: %q", p.Keys)
	}

	p.AddKey(Jade)
	if p.Keys.String() != Copper.String()+"|"+Jade.String() {
		t.Fatalf("+jade: %q", p.Keys)
	}

	p.RemoveKey(Copper)
	if p.Keys.String() != Jade.String() {
		t.Fatalf("-copper: %q", p.Keys)
	}

	if p.HasKey(Copper) {
		t.Fatalf("copper in %q", p.Keys)
	}

	if !p.HasKey(Jade) {
		t.Fatalf("jade not in %q", p.Keys)
	}
}
