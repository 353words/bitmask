package main

// DocumentPermissions are permissions set on a document
type DocumentPermissions uint8

// Available permissions
const (
	Locked DocumentPermissions = 1 << iota
	GroupReadable
	GroupWritable
	AllReadable
	AllWritable
)

func (p *DocumentPermissions) Set(perm DocumentPermissions) {
	*p = *p | perm
}

func (p *DocumentPermissions) Clear(perm DocumentPermissions) {
	*p = *p & (^perm)
}

func (p DocumentPermissions) IsSet(perm DocumentPermissions) bool {
	return p&perm != 0
}
