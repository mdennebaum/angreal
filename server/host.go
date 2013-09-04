package server

import (
	"github.com/mdennebaum/angreal/util"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Host struct {
	conf         *util.DynMap
	url          string
	port         string
	root         string
	static       string
	headers      *util.DynMap
	proxies      []*httputil.ReverseProxy
	backends     *util.DynMap
	proxyChannel chan *httputil.ReverseProxy
}

func NewHost(config *util.DynMap) *Host {
	h := new(Host)
	h.conf = config
	h.proxyChannel = make(chan *httputil.ReverseProxy)
	return h
}

func (this *Host) Init() {
	this.url, _ = this.conf.GetString("url")
	this.port, _ = this.conf.GetString("port")
	this.root, _ = this.conf.GetString("root")
	this.static, _ = this.conf.GetString("static")
	this.backends, _ = this.conf.GetDynMap("backends")
	this.headers, _ = this.conf.GetDynMap("headers")
	this.initLog()
	this.initBackends()
	this.initHostHandler()
}

func (this *Host) initLog() {
	// if accessLog, ok := this.conf.GetString("access_log"); ok {

	// }
	// if errorLog, ok := this.conf.GetString("error_log"); ok {

	// }
}

func (this *Host) addHeaders(w http.ResponseWriter) http.ResponseWriter {
	if this.headers != nil {
		for k, v := range this.headers.Map {
			w.Header().Set(k, v.(string))
		}
	}
	return w
}

//configure our reverse proxy backends
func (this *Host) initBackends() {

	//check if we have any backends to proxy requests to
	if this.backends != nil {

		//loop over configs
		for _, v := range this.backends.Map {
			//grab the server url and parse it
			serverUrl, err := url.Parse(v.(string))

			//if its not a valid url log the error
			if err != nil {
				log.Println(err)
			} else {
				//create a new proxy for our backend and add it to our proxy slice
				this.proxies = append(this.proxies, httputil.NewSingleHostReverseProxy(serverUrl))
			}
		}

		//loop over proxies and block till one is needed
		go func() {
			for {
				for _, p := range this.proxies {
					this.proxyChannel <- p
				}
			}
		}()
	}
}

//get the next available proxy from our round robin channel
func (this *Host) NextProxy() *httputil.ReverseProxy {
	return <-this.proxyChannel
}

// Provide proxying of a url. Reverse proxy just masks the path
func (this *Host) initHostHandler() {
	http.HandleFunc(this.url+":"+this.port+"/", func(w http.ResponseWriter, r *http.Request) {
		w = this.addHeaders(w)
		if this.static != "" {
			if strings.HasPrefix(r.URL.Path, this.static) {
				http.ServeFile(w, r, this.root+r.URL.Path)
				return
			}
		}
		//get next proxy
		rewriteProxy := this.NextProxy()
		r.RequestURI = strings.Replace(r.RequestURI, "/", "", 1)
		r.URL.Path = strings.Replace(r.URL.Path, "/", "", 1)
		rewriteProxy.ServeHTTP(w, r)
	})
}
