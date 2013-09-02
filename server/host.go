package server

/*
TODO this should be a obj that encapsulates all the logic for a host
logging, file server setup, reverse proxy setup, header addition
*/

import (
	"github.com/mdennebaum/angreal/util"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Host struct {
	conf   *util.DynMap
	url    string
	port   string
	root   string
	static string
}

func NewHost(config *util.DynMap) *Host {
	h := Host{config, "", "", "", ""}
	return &h
}

func (this *Host) Init() {
	this.url, _ = this.conf.GetString("url")
	this.port, _ = this.conf.GetString("port")
	this.root, _ = this.conf.GetString("root")
	this.static, _ = this.conf.GetString("static")

	this.initLog()
	this.initStatic()
	this.initBackends()

}

func (this *Host) initLog() {
	// if accessLog, ok := this.conf.GetString("access_log"); ok {

	// }
	// if errorLog, ok := this.conf.GetString("error_log"); ok {

	// }
}

func (this *Host) addHeaders(w http.ResponseWriter) {

}

func (this *Host) initBackends() {

	serverUrl, _ := url.Parse("http://test1.localhost.com:8000")
	this.addProxy("test1.localhost.com:8080/", serverUrl, true)
}

// Provide proxying of a url. Reverse proxy just masks the path
func (this *Host) addProxy(path string, url *url.URL, reverse bool) {

	rewriteProxy := httputil.NewSingleHostReverseProxy(url)
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		log.Println("in prox handler")
		if reverse == true {
			r.RequestURI = strings.Replace(r.RequestURI, path, "", 1)
			r.URL.Path = strings.Replace(r.URL.Path, path, "", 1)
		}
		rewriteProxy.ServeHTTP(w, r)
	})
}

func (this *Host) initStatic() {
	http.HandleFunc(this.url+":"+this.port+"/"+this.static+"/", this.getStaticHandler())
}

func (this *Host) getStaticHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		this.addHeaders(w)
		// this.log.Access(r)
		http.ServeFile(w, r, this.root+"/"+this.static+"/"+r.URL.Path[1:])
	}
}
