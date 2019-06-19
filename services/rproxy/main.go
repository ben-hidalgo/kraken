// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxy(target string) *Prox {
	url, _ := url.Parse(target)
	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")
	// p.proxy.Transport = &myTransport{}
	p.proxy.ServeHTTP(w, r)
}

func main() {
	const (
		defaultPort        = ":9090"
		defaultPortUsage   = "default server port, ':9090'"
		defaultTarget      = "http://127.0.0.1:9091"
		defaultTargetUsage = "default redirect url, 'http://127.0.0.1:9091'"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	redirecturl := flag.String("url", defaultTarget, defaultTargetUsage)

	flag.Parse()

	fmt.Println("server will run on :", *port)
	fmt.Println("redirecting to :", *redirecturl)

	// proxy
	proxy := NewProxy(*redirecturl)

	http.HandleFunc("/scenes", proxy.handle)
	// http.HandleFunc("/", proxy.handle)

	log.Fatal(http.ListenAndServe(*port, nil))
}
