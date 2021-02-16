package main

import (
	"github.com/jvitoroc/todo-go/server"
)

func main() {
	server := server.NewServer()
	server.Start()
}
