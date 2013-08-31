package server

import (
	"github.com/mdennebaum/angreal/util"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type Server struct {
	hosts  []*Host
	Config util.DynMap
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

	//try and get the global config
	if gc, ok := this.Config["global"]; ok {

		//cast it to a map
		global := gc.(map[string]interface{})
		//check if we have an explicit processor setting
		if procs, ok := global["procs"]; ok {
			//conv to int
			count, _ := strconv.Atoi(procs.(string))
			//use the explicit processor count
			runtime.GOMAXPROCS(count)
			return
		}
	}

	//use all avail processors
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (this *Server) Listen() {

	//  if err := http.ListenAndServeTLS(f.CertFile, f.KeyFile); err != nil {
	//    log.Printf("Starting HTTPS frontend %s failed: %v", f.Name, err)
	//  }

	log.Println("started on port 8000")
	srv := &http.Server{
		Addr:        ":8000",
		ReadTimeout: 30 * time.Second,
	}
	srv.ListenAndServe()

}

func (this *Server) setupHosts() {
	//loop over config vhosts and call setupHost for each
	if h, ok := this.Config["hosts"]; ok {

		//cast it to a map
		hosts := h.([]interface{})
		for _, host_config := range hosts {
			host := NewHost(host_config.(map[string]interface{}))
			host.Init()
			this.hosts = append(this.hosts, host)
		}
	}
}

func (this *Server) loadConfig() {
	config := util.NewConfig("./angreal.conf")
	config.Load()
	this.Config = config.Data
}
