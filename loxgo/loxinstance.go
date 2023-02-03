package main

type LoxInstance struct {
	*LoxClass
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class}
}

//	func (li *LoxInstance) String() string {
//		return li.name + " instance"
//	}
func (li LoxInstance) String() string {
	return li.name + " instance"
}
