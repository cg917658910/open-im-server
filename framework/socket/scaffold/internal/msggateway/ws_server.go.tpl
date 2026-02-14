package msggateway

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

type Server struct {
	upgrader   websocket.Upgrader
	register   chan *Client
	unregister chan *Client
	kick       chan *Client
	users      map[string]map[*Client]struct{}
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		register: make(chan *Client, 1024),
		unregister: make(chan *Client, 1024),
		kick: make(chan *Client, 1024),
		users: make(map[string]map[*Client]struct{}),
	}
}

func (s *Server) Run() {
	for {
		select {
		case c := <-s.register:
			if _, ok := s.users[c.UserID]; !ok {
				s.users[c.UserID] = map[*Client]struct{}{}
			}
			s.users[c.UserID][c] = struct{}{}
		case c := <-s.unregister:
			if m, ok := s.users[c.UserID]; ok {
				delete(m, c)
				if len(m) == 0 {
					delete(s.users, c.UserID)
				}
			}
			_ = c.Conn.Close()
		case c := <-s.kick:
			_ = c.Conn.Close()
		}
	}
}
