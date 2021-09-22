package peer

type Session interface {
	Send(msg interface{}) error
}
