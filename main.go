package main

import (
	"github.com/mdennebaum/angreal/server"
)

func main() {
	angreal := new(server.Server)
	angreal.Init()
	angreal.Listen()
}
