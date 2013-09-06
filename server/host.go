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
	url          string
	root         string
	static       string
	port         string
	reqLog       *util.Log
	conf         *util.DynMap
	backends     *util.DynMap
	headers      *util.DynMap
	proxies      []*httputil.ReverseProxy
	proxyChannel chan *httputil.ReverseProxy
}

func NewHost(config *util.DynMap, port string) *Host {
	h := new(Host)
	h.conf = config
	h.port = port
	h.reqLog = util.NewLog()
	h.proxyChannel = make(chan *httputil.ReverseProxy)
	return h
}

func (h *Host) Init() {
	h.url = h.conf.MustString("url", "/")
	h.root, _ = h.conf.GetString("root")
	h.static, _ = h.conf.GetString("static")
	h.backends, _ = h.conf.GetDynMap("backends")
	h.headers, _ = h.conf.GetDynMap("headers")
	h.initLog()
	h.initBackends()
	h.initHostHandler()
}

func (h *Host) initLog() {
	if accessLog, ok := h.conf.GetString("access_log"); ok {
		h.reqLog.AddLog("access", accessLog)
	}
	if errorLog, ok := h.conf.GetString("error_log"); ok {
		h.reqLog.AddLog("error", errorLog)
	}
}

func (h *Host) addHeaders(w http.ResponseWriter) http.ResponseWriter {
	if h.headers != nil {
		for k, v := range h.headers.Map {
			w.Header().Set(k, v.(string))
		}
	}
	return w
}

//configure our reverse proxy backends
func (h *Host) initBackends() {

	//check if we have any backends to proxy requests to
	if h.backends != nil {

		//loop over configs
		for _, v := range h.backends.Map {
			//grab the server url and parse it
			serverUrl, err := url.Parse(v.(string))

			//if its not a valid url log the error
			if err != nil {
				log.Println(err)
			} else {
				//create a new proxy for our backend and add it to our proxy slice
				h.proxies = append(h.proxies, httputil.NewSingleHostReverseProxy(serverUrl))
			}
		}

		//loop over proxies and block till one is needed
		go func() {
			for {
				for _, p := range h.proxies {
					h.proxyChannel <- p
				}
			}
		}()
	}
}

//get the next available proxy from our round robin channel
func (h *Host) NextProxy() *httputil.ReverseProxy {
	return <-h.proxyChannel
}

// Provide proxying of a url. Reverse proxy just masks the path
func (h *Host) initHostHandler() {
	http.HandleFunc(h.url+":"+h.port+"/", func(w http.ResponseWriter, r *http.Request) {
		h.reqLog.Access(r)
		w = h.addHeaders(w)
		if h.static != "" {
			if strings.HasPrefix(r.URL.Path, h.static) {
				http.ServeFile(w, r, h.root+r.URL.Path)
				return
			}
		}
		//get next proxy
		rewriteProxy := h.NextProxy()
		r.RequestURI = strings.Replace(r.RequestURI, "/", "", 1)
		r.URL.Path = strings.Replace(r.URL.Path, "/", "", 1)
		rewriteProxy.ServeHTTP(w, r)
	})
}
