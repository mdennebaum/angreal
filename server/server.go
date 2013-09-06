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

func (s *Server) Init() {
	s.loadConfig()
	s.initProcs()
	s.setupHosts()
}

//init processors
func (s *Server) initProcs() {
	//use whats in the conf or all avail processors
	runtime.GOMAXPROCS(s.conf.MustInt("global.procs", runtime.NumCPU()))
}

func (s *Server) Listen() {
	//TODO serve https
	//  if err := http.ListenAndServeTLS(f.CertFile, f.KeyFile); err != nil {
	//    log.Printf("Starting HTTPS frontend %s failed: %v", f.Name, err)
	//  }

	//TODO get the ports for this server
	log.Println("started on port 8080")
	srv := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 30 * time.Second,
	}
	srv.ListenAndServe()

}

func (s *Server) setupHosts() {
	//loop over config hosts and setup new host for each
	if hosts, ok := s.conf.GetDynMapSlice("hosts"); ok {
		for _, host := range hosts {
			h := NewHost(host)
			h.Init()
			s.hosts = append(s.hosts, h)
		}
	}
}

func (s *Server) loadConfig() {
	s.conf = util.NewConfig("./angreal.conf")
	s.conf.Load()
}
