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
	conf          *util.DynMap
	url           string
	port          string
	root          string
	static        string
	proxies       []*httputil.ReverseProxy
	backends      *util.DynMap
	proxyChannel  chan *httputil.ReverseProxy
}

func NewHost(config *util.DynMap) *Host {
	h := new(Host)
	h.conf = config
	return h
}

func (this *Host) Init() {
	this.url, _ = this.conf.GetString("url")
	this.port, _ = this.conf.GetString("port")
	this.root, _ = this.conf.GetString("root")
	this.static, _ = this.conf.GetString("static")
	this.backends, _ = this.conf.GetDynMap("backends")
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

func (this *Host) addHeaders(w http.ResponseWriter) {

}

func (this *Host) initBackends() {
	if this.backends != nil {
		for _, v := range this.backends.Map {
			serverUrl, err := url.Parse(v.(string))
			if err != nil {
				log.Println(err)
			} else {
				this.proxies = append(this.proxies, httputil.NewSingleHostReverseProxy(serverUrl))
			}
		}
	}
	
	go func(){
		for {
		  for p in this.proxies {
		    this.proxyChannel <- p
		  }
		}
	}
}

func (this *Host) Next() *httputil.ReverseProxy{
  return <- this.proxyChannel
}

//grab the next round robin backend to handle req
func (this *Host) getProxyBackend() *httputil.ReverseProxy {
	return this.proxies[this.proxyPosition.Next()]
}

// Provide proxying of a url. Reverse proxy just masks the path
func (this *Host) initHostHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if this.static != "" {
			if strings.HasPrefix(r.URL.Path, this.static) {
				http.ServeFile(w, r, this.root+r.URL.Path)
				return
			}
		}
		//get next proxy
		rewriteProxy := this.Next()
		r.RequestURI = strings.Replace(r.RequestURI, "/", "", 1)
		r.URL.Path = strings.Replace(r.URL.Path, "/", "", 1)
		rewriteProxy.ServeHTTP(w, r)
	})
}
