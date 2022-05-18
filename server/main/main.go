package main

import (
	"chatroom/server"
)

func main() {
	var s server.ChatServer
	s = server.NewServer()
	s.Listen(":3333")
	//start server
	s.Start()
}
