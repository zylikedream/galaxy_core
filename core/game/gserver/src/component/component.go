package component

type Component interface {
	Name() string
	Load()
	Unload()
}

type IPersit interface {
	Table() string
	Serialize(interface{}) interface{}
	Unserialize(interface{})
}
