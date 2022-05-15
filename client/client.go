package client

import "chatroom/tube"

type ChatClient interface {
	Dial(address string) error //Dial 建立連線並且建立通訊協議的reader和writer
	Start()
	Close()
	Send(command interface{}) error
	SetName(name string) error
	SendMessage(message string) error
	Error() chan error
	Incoming() chan tube.MessageCommand
}
