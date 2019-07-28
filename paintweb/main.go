package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/qiniu/qpaint/paintdom"
)

func newReverseProxy(baseURL string) *httputil.ReverseProxy {
	rpURL, _ := url.Parse(baseURL)
	return httputil.NewSingleHostReverseProxy(rpURL)
}

func handleDefault(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.ServeFile(w, req, "www/index.htm")
		return
	}
	req.URL.RawQuery = "" // skip "?params"
	wwwServer.ServeHTTP(w, req)
}

var (
	apiReverseProxy = newReverseProxy("http://localhost:9999")
	wwwServer       = http.FileServer(http.Dir("www"))
)

func main() {
	go paintdom.Main()
	http.Handle("/api/", http.StripPrefix("/api/", apiReverseProxy))
	http.HandleFunc("/", handleDefault)
	http.ListenAndServe(":8888", nil)
}
