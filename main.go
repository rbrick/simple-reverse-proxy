package main

import "net/http"

func init() {

}

func main() {
	http.ListenAndServe(":80", NewProxiedHandler("hosts.json"))
}
