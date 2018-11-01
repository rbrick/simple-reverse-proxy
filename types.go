package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxiedHandler struct {
	ReverseProxies map[string]*httputil.ReverseProxy
}

func (p *ProxiedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if v, ok := p.ReverseProxies[req.Host]; ok {
		v.ServeHTTP(w, req)
	}
}

func NewProxiedHandler(hostFile string) *ProxiedHandler {
	data, err := ioutil.ReadFile(hostFile)

	if err != nil {
		log.Fatalln(err)
	}

	var hostRedirectMap map[string]*url.URL

	json.Unmarshal(data, &hostRedirectMap)

	reversedProxies := map[string]*httputil.ReverseProxy{}

	for k, v := range hostRedirectMap {
		reversedProxies[k] = httputil.NewSingleHostReverseProxy(v)
	}

	return &ProxiedHandler{
		ReverseProxies: reversedProxies,
	}
}
