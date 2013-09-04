package main

import (
	"github.com/mdennebaum/angreal/server"
)

//TODO use flag to pass in startup options... config
func main() {
	angreal := server.NewServer()
	angreal.Init()
	angreal.Listen()
}
