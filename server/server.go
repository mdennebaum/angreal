package server

import (
	"github.com/mdennebaum/angreal/util"
	"log"
	"net/http"
	"runtime"
	"time"
)

type Server struct {
	hosts []*Host
	conf  *util.Config
	log   *util.Log
}

func NewServer() *Server {
	s := Server{}
	return &s
}

func (this *Server) Init() {
	this.loadConfig()
	this.initProcs()
	this.setupHosts()
}

//init processors
func (this *Server) initProcs() {

	//check if we have an explicit processor setting
	if procs, ok := this.conf.GetInt("global.procs"); ok {
		//use the explicit processor count
		runtime.GOMAXPROCS(procs)
		return
	}

	//use all avail processors
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func (this *Server) Listen() {
	//TODO serve https
	//  if err := http.ListenAndServeTLS(f.CertFile, f.KeyFile); err != nil {
	//    log.Printf("Starting HTTPS frontend %s failed: %v", f.Name, err)
	//  }

	log.Println("started on port 8080")
	srv := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 30 * time.Second,
	}
	srv.ListenAndServe()

}

func (this *Server) setupHosts() {

	//loop over config hosts and setup new host for each
	if hosts, ok := this.conf.GetDynMapSlice("hosts"); ok {
		for _, host := range hosts {
			h := NewHost(host)
			h.Init()
			this.hosts = append(this.hosts, h)
		}
	}
}

func (this *Server) loadConfig() {
	this.conf = util.NewConfig("./angreal.conf")
	this.conf.Load()
}
