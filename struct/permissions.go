package main

import (
	"fmt"
	"unsafe"
)

type DocumentPermissions struct {
	Locked        bool
	GroupReadable bool
	GroupWritable bool
	AnyReadable   bool
	AnyWritable   bool
}

func main() {
	perms := DocumentPermissions{}
	fmt.Println(unsafe.Sizeof(perms))
}
