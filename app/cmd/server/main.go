package main

import "github.com/redis-go/app/internal/server"

func main() {
	srv := server.New(":6379")
	srv.Start()
}
