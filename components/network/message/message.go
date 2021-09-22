package message

type MessageCodec interface {
	Decode(ID int, data []byte) (interface{}, error)
	Encode(packet interface{}) ([]byte, int, error)
	ReisterPacket(ID int, packet interface{}) error
}

type Message struct {
	ID      int
	Type    int
	Payload []byte
	Msg     interface{}
}
