package main

import (
	"testing"
)

func isIn(p DocumentPermissions, perms []DocumentPermissions) bool {
	for _, op := range perms {
		if p == op {
			return true
		}
	}
	return false
}

func checkPerms(t *testing.T, perm DocumentPermissions, set []DocumentPermissions) {
	for bit := 0; bit < 7; bit++ {
		p := DocumentPermissions(1 << bit)
		if isIn(p, set) {
			if !perm.IsSet(p) {
				t.Fatalf("%v: %v not set", perm, p)
			}
		} else {
			if perm.IsSet(p) {
				t.Fatalf("%v: %v is set", perm, p)
			}
		}
	}
}

func TestPermissions(t *testing.T) {
	perm := Locked | AllReadable
	checkPerms(t, perm, []DocumentPermissions{Locked, AllReadable})

	perm.Set(GroupWritable)
	checkPerms(t, perm, []DocumentPermissions{Locked, AllReadable, GroupWritable})

	perm.Clear(Locked)
	checkPerms(t, perm, []DocumentPermissions{AllReadable, GroupWritable})
}
