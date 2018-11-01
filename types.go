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
	if v, ok := p.ReverseProxies[req.URL.Host]; ok {
		v.ServeHTTP(w, req)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(http.StatusText(http.StatusNotFound)))
}

func clean(s []string) []string {
	var x []string
	for _, str := range s {
		if str != "" {
			x = append(x, str)
		}
	}
	return x
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
		realReverseProxy := httputil.NewSingleHostReverseProxy(v)

		loggingProxy := &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				realReverseProxy.Director(r) // run our director first
				if r.TLS != nil && v.Scheme == "http" {
					// panic!
					panic("http over https")
				} else {
					log.Printf("sending proxy request on %s://%s%s to %s\n", r.URL.Scheme, k, r.URL.Path, v.Host)
				}
			},
		}

		reversedProxies[k] = loggingProxy
	}

	return &ProxiedHandler{
		ReverseProxies: reversedProxies,
	}
}
