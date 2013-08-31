package server

/*
TODO this should be a obj that encapsulates all the logic for a host
logging, file server setup, reverse proxy setup, header addition
*/

// import (
// 	"log"
// 	"net/http"
// 	"runtime"
// 	"strconv"
// 	"time"
// )

// type Server struct {
// 	Config map[string]interface{}
// }

// func (this *Server) Init() {
// 	this.loadConfig()
// 	this.initProcs()
// 	this.setupHosts()
// }

// //init processors
// func (this *Server) initProcs() {

// 	//try and get the global config
// 	if gc, ok := this.Config["global"]; ok {

// 		//cast it to a map
// 		global := gc.(map[string]interface{})
// 		//check if we have an explicit processor setting
// 		if procs, ok := global["procs"]; ok {
// 			//conv to int
// 			count, _ := strconv.Atoi(procs.(string))
// 			//use the explicit processor count
// 			runtime.GOMAXPROCS(count)
// 			return
// 		}
// 	}

// 	//use all avail processors
// 	runtime.GOMAXPROCS(runtime.NumCPU())
// }

// func (this *Server) Listen() {

// 	//  if err := http.ListenAndServeTLS(f.CertFile, f.KeyFile); err != nil {
// 	//    log.Printf("Starting HTTPS frontend %s failed: %v", f.Name, err)
// 	//  }

// 	log.Println("started on port 8000")
// 	srv := &http.Server{
// 		Addr:        ":8000",
// 		ReadTimeout: 30 * time.Second,
// 	}
// 	srv.ListenAndServe()

// }

// func (this *Server) setupHosts() {
// 	//loop over config vhosts and call setupHost for each
// 	if vh, ok := this.Config["vhosts"]; ok {

// 		//cast it to a map
// 		vhosts := vh.([]interface{})
// 		for _, host := range vhosts {
// 			this.setupHost(host)
// 		}
// 	}
// }

// func (this *Server) getStaticHandler(root string) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, root+"/"+r.URL.Path[1:])
// 	}
// }

// func (this *Server) setupHost(host interface{}) {
// 	conf := host.(map[string]interface{})
// 	url := conf["url"].(string)
// 	//static := conf["static"].(string)
// 	port := conf["port"].(string)
// 	root := conf["root"].(string)

// 	http.HandleFunc(url+":"+port+"/", this.getStaticHandler(root))
// 	// http.HandleFunc(url+":"+port+"/", func(w http.ResponseWriter, r *http.Request) {
// 	//  fmt.Fprintf(w, url)
// 	// })
// }

// func (this *Server) loadConfig() {
// 	config := NewConfig("./raptrix.conf")
// 	config.Load()
// 	this.Config = config.Data
// }
