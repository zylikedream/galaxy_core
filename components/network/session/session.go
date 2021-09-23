package session

type Session interface {
	Send(msg interface{}) error
	Close()
}
