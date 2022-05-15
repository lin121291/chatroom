package server

type ChatServer interface {
	Listen(address string) error
	Start()
	Broadcast(command interface{}) error
	Close()
}
