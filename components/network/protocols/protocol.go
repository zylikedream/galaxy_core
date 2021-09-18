package protocols

type Protocols struct {
	id2proto map[int]Proto
}

type Proto struct {
	Name string
	Type interface{}
}
