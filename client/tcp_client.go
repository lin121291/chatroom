package client

import (
	"chatroom/tube"
	"log"
	"net"
)

type TcpChatClient struct {
	conn      net.Conn
	cmdReader *tube.CommandReader
	cmdWriter *tube.CommandWriter
	name      string
	error     chan error
	incoming  chan tube.MessageCommand
}

func NewClient() *TcpChatClient {
	return &TcpChatClient{
		incoming: make(chan tube.MessageCommand),
		error:    make(chan error),
	}
}

func (c *TcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err == nil {
		c.conn = conn
		c.cmdReader = tube.NewCommandReader(conn)
		c.cmdWriter = tube.NewCommandWriter(conn)
	}

	return err
}

func (c *TcpChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()

		if err != nil {
			c.error <- err
			break
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case tube.MessageCommand:
				c.incoming <- v
			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}

func (c *TcpChatClient) Close() {
	c.conn.Close()
}

func (c *TcpChatClient) Incoming() chan tube.MessageCommand {
	return c.incoming
}

func (c *TcpChatClient) Error() chan error {
	return c.error
}

func (c *TcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *TcpChatClient) SetName(name string) error {
	return c.Send(tube.NameCommand{name})
}

func (c *TcpChatClient) SendMessage(message string) error {
	return c.Send(tube.SendCommand{
		Message: message,
	})
}
