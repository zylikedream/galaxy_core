package codec

type PacketDecoder interface {
	Decode(ID int, data []byte) (interface{}, error)
}

type PacketEncoder interface {
	Encode(msg interface{}) ([]byte, error)
}
