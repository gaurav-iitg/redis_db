package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/redis-go/app/internal/dispatcher"
	"github.com/redis-go/app/internal/resp"
)

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(c)
	for {
		cmd, err := resp.ReadRESP(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Error reading request:", err)
			return
		}

		response, err := s.dispatcher.Execute(cmd)
		if err != nil {
			fmt.Println("Error handling response:", err)
			return
		}

		payload, err := resp.EncodeResp(response)
		if err != nil {
			fmt.Println("Error encoding response:", err)
			return
		}

		if _, err := c.Write(payload); err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
	}
}

type Server struct {
	addr       string
	dispatcher dispatcher.Dispatcher
}

func New(addr string) *Server {
	dispatcher := dispatcher.New()
	return &Server{addr: addr, dispatcher: *dispatcher}
}

func (s *Server) Start() {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Printf("Failed to bind to port %s: %v\n", s.addr, err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			os.Exit(1)
		}
		go s.handleConn(conn)
	}
}
