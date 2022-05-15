package server

import (
	"chatroom/tube"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type client struct {
	conn   net.Conn
	name   string
	writer *tube.CommandWriter
}

type TcpChatServer struct {
	listener net.Listener //監聽port
	clients  []*client
	mutex    *sync.Mutex
}

var (
	UnknownClient = errors.New("Unknown client")
)

func NewServer() *TcpChatServer {
	return &TcpChatServer{
		mutex: &sync.Mutex{},
	}
}

//監聽port
func (s *TcpChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)

	if err == nil {
		s.listener = l
	}

	log.Printf("Listening on %v", address)

	return err
}

func (s *TcpChatServer) Close() {
	s.listener.Close()
}

func (s *TcpChatServer) Start() {
	//在這裡做無限迴圈等待
	for {
		// XXX: need a way to break the loop
		conn, err := s.listener.Accept()

		if err != nil {
			log.Print(err)
		} else {
			// handle connection
			client := s.accept(conn) //設定本次client
			go s.serve(client)       //在這裡開一個新的
		}
	}
}

//在這裡還不知道怎麼推送到每個client
func (s *TcpChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		// TODO: handle error here?
		client.writer.Write(command)
	}

	return nil
}

//建立對應的client跟蹤使用者
func (s *TcpChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	//為什麼這邊會需要mutex
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//一開始只有設定連線和寫入
	client := &client{
		conn:   conn,
		writer: tube.NewCommandWriter(conn), //這邊net.conn轉成了io.write
	}

	s.clients = append(s.clients, client) //把這次的client記起來

	return client
}

func (s *TcpChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// remove the connections from clients array
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

func (s *TcpChatServer) serve(client *client) {

	//創建NewCommandReader並連線上
	cmdReader := tube.NewCommandReader(client.conn) //net.Conn 也是一個 Reader

	defer s.remove(client)

	for {
		//client先write過來了東西，所以先read
		cmd, err := cmdReader.Read()

		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) { //cmd是哪一種struct
			case tube.SendCommand:
				go s.Broadcast(tube.MessageCommand{ //傳送訊息給全部人
					Message: v.Message,
					Name:    client.name,
				})

			//在新開的thread上寫一個新的名子
			case tube.NameCommand: //設定名子
				client.name = v.Name
			}
		}

		if err == io.EOF {
			break
		}
	}
}
