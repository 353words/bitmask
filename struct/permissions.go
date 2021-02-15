package main

import (
	"fmt"
	"unsafe"
)

// DocumentPermissions are permissions on a document.
type DocumentPermissions struct {
	Locked        bool
	GroupReadable bool
	GroupWritable bool
	AnyReadable   bool
	AnyWritable   bool
}

func main() {
	var perms DocumentPermissions
	fmt.Println(unsafe.Sizeof(perms))
}
