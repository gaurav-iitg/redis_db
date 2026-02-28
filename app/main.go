package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal/dispatcher"
	"github.com/codecrafters-io/redis-starter-go/app/internal/resp"
)

func handleConn(c net.Conn) {
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

		response, err := dispatcher.Execute(cmd)
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

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Printf("Failed to bind to port 6379: %v\n", err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			os.Exit(1)
		}
		go handleConn(conn)
	}
}
