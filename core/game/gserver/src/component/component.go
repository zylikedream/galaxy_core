package component

type Component interface {
	Name() string
	Load() error
	Unload() error
}

type IPersit interface {
	Table() string
	Serialize(interface{}) interface{}
	Unserialize(interface{})
}
