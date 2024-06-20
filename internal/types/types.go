package types

import "github.com/scylladb/go-set/strset"

type File struct {
	ImportPathSet *strset.Set
	Interfaces    []*Interface
}

type Interface struct {
	Name    string
	Methods []*Method
	// TODO: embed interface
	ImportPathSet *strset.Set
}

type Method struct {
	Name    string
	Args    []*Tuple
	Results []*Tuple
}

type Tuple struct {
	Name string
	Type string
}
