package server

/*
TODO this should be a obj that encapsulates all the logic for a host
logging, file server setup, reverse proxy setup, header addition
*/

import (
	// "log"
	"net/http"
)

type Host struct {
	Config map[string]interface{}
}

func NewHost(config map[string]interface{}) *Host {
	h := Host{config}
	return &h
}

func (this *Host) Init() {
	this.initLog()
	this.initBackends()
	this.initStatic()
}

func (this *Host) initLog() {}

func (this *Host) addHeaders() {

}

func (this *Host) initBackends() {}

func (this *Host) initStatic() {
	url := this.Config["url"].(string)
	port := this.Config["port"].(string)
	root := this.Config["root"].(string)
	http.HandleFunc(url+":"+port+"/", this.getStaticHandler(root))
}

func (this *Host) getStaticHandler(root string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		this.addHeaders()
		http.ServeFile(w, r, root+"/"+r.URL.Path[1:])
	}
}
