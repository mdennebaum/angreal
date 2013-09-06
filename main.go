package main

import (
	"flag"
	"github.com/mdennebaum/angreal/server"
)

//the configuration file path var
var configPath string

//pass in startup option for config path
func init() {
	flag.StringVar(&configPath, "config", "./angreal.conf", "path to your angreal conf file")
}

func main() {

	//parse the command line args
	flag.Parse()

	//get a new server object
	angreal := server.NewServer(configPath)

	//init and start to listen for connections
	(angreal.Init()).Listen()
}
