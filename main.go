package main

import (
	"github.com/mdennebaum/angreal/server"
)

func main() {
	angreal := server.NewServer()
	angreal.Init()
	angreal.Listen()
}
