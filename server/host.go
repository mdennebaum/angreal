package server

/*
TODO this should be a obj that encapsulates all the logic for a host
logging, file server setup, reverse proxy setup, header addition
*/

import (
	"github.com/mdennebaum/angreal/util"
	"net/http"
)

type Host struct {
	conf *util.DynMap
}

func NewHost(config *util.DynMap) *Host {
	h := Host{config}
	return &h
}

func (this *Host) Init() {
	this.initLog()
	this.initBackends()
	this.initStatic()
}

func (this *Host) initLog() {
	if accessLog, ok := this.conf.GetString("access_log"); ok {

	}

	if errorLog, ok := this.conf.GetString("error_log"); ok {

	}
}

func (this *Host) addHeaders() {

}

func (this *Host) initBackends() {

}

func (this *Host) initStatic() {
	url, _ := this.conf.GetString("url")
	port, _ := this.conf.GetString("port")
	root, _ := this.conf.GetString("root")
	http.HandleFunc(url+":"+port+"/", this.getStaticHandler(root))
}

func (this *Host) getStaticHandler(root string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		this.addHeaders()
		http.ServeFile(w, r, root+"/"+r.URL.Path[1:])
	}
}
