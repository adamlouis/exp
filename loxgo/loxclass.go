package main

type LoxClass struct {
	name string
}

func (lc *LoxClass) String() string {
	return lc.name
}
